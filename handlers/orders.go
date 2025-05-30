package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ringover_assignment/db"
	"github.com/Ringover_assignment/models"
	"github.com/gorilla/mux"
)

func HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	err := cancelOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resError, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(resError)
		return
	}
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]string{"message": "Order canceled"})
	w.Write(res)
}

func HandleGetOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	book, err := getOrderBook(symbol)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		resError, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.Write(resError)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func HandleGetTrades(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	trades, err := getTrades(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(trades)
}

func HandleGetOrderStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	order, err := getOrderStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func HandlePlaceOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// Validate required fields (basic validation)
	if order.Symbol == "" || order.Side == "" || order.Type == "" || order.Quantity <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	
	res, err := processOrder(order)
	if err != nil {
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Order received"})
	json.NewEncoder(w).Encode(res)
}

func processOrder(order models.Order) (*models.Order, error) {
	order.RemainingQty = order.Quantity
	order.Status = "open"

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`
		INSERT INTO orders (symbol, side, type, price, quantity, remaining_quantity, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		order.Symbol, order.Side, order.Type, order.Price, order.Quantity, order.RemainingQty, order.Status)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}

	order.ID, _ = res.LastInsertId()
	if err := matchOrders(tx, &order); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}
	

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func matchOrders(tx *sql.Tx, order *models.Order) error {
	var query string
	var args []interface{}

	if order.Type == models.MarketOrder {
		if order.Side == models.BuySide {
			query = `SELECT id, price, remaining_quantity FROM orders
				WHERE symbol=? AND side='sell' AND status='open'
				ORDER BY price ASC, created_at ASC`
		} else {
			query = `SELECT id, price, remaining_quantity FROM orders
				WHERE symbol=? AND side='buy' AND status='open'
				ORDER BY price DESC, created_at ASC`
		}
		args = []interface{}{order.Symbol}
	} else {
		if order.Side == models.BuySide {
			query = `SELECT id, price, remaining_quantity FROM orders
				WHERE symbol=? AND side='sell' AND status='open' AND price <= ?
				ORDER BY price ASC, created_at ASC`
			args = []interface{}{order.Symbol, order.Price}
		} else {
			query = `SELECT id, price, remaining_quantity FROM orders
				WHERE symbol=? AND side='buy' AND status='open' AND price >= ?
				ORDER BY price DESC, created_at ASC`
			args = []interface{}{order.Symbol, order.Price}
		}
	}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	type Match struct {
		ID    int64
		Price float64
		Qty   int
	}
	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.Price, &m.Qty); err != nil {
			return err
		}
		matches = append(matches, m)
	}

	var tradePlaceholders []string
	var tradeValues []interface{}

	var fullyFilledIDs []int64
	var partiallyFilled Match
	var partialQtyUsed int
	partialMatched := false

	for _, match := range matches {
		if order.RemainingQty <= 0 {
			break
		}
		matchQty := min(order.RemainingQty, match.Qty)

		buyID := chooseBuyID(order.Side, order.ID, match.ID)
		sellID := chooseSellID(order.Side, order.ID, match.ID)

		tradePlaceholders = append(tradePlaceholders, "(?, ?, ?, ?, ?)")
		tradeValues = append(tradeValues, buyID, sellID, order.Symbol, match.Price, matchQty)

		if matchQty == match.Qty {
			fullyFilledIDs = append(fullyFilledIDs, match.ID)
		} else {
			partialMatched = true
			partiallyFilled = match
			partialQtyUsed = matchQty
		}

		order.RemainingQty -= matchQty
	}

	
	if len(tradePlaceholders) > 0 {
		insertQuery := fmt.Sprintf(`INSERT INTO trades (buy_order_id, sell_order_id, symbol, price, quantity) VALUES %s`,
			strings.Join(tradePlaceholders, ","))
		if _, err := tx.Exec(insertQuery, tradeValues...); err != nil {
			return err
		}
	}

	
	if len(fullyFilledIDs) > 0 {
		placeholder := strings.Repeat("?,", len(fullyFilledIDs))
		placeholder = placeholder[:len(placeholder)-1] // remove trailing comma
		args := make([]any, len(fullyFilledIDs)*2)
		for i, id := range fullyFilledIDs {
			args[i] = 0                      
			args[i+len(fullyFilledIDs)] = id 
		}
		arrayId := make([]any, len(fullyFilledIDs))
		for i, id := range fullyFilledIDs {
			arrayId[i] = id
		}

		// `remaining_quantity - 0 = 0` is true for already matched qty == full qty
		updateQuery := fmt.Sprintf(`
			UPDATE orders SET remaining_quantity = 0, status = 'filled'
			WHERE id IN (%s)	
		`, placeholder)
		if _, err := tx.Exec(updateQuery, arrayId...); err != nil {
			return err
		}
	}

	// Update one partially matched order if exists
	if partialMatched {
		_, err := tx.Exec(`
			UPDATE orders
			SET remaining_quantity = remaining_quantity - ?, status = 'open'
			WHERE id = ?`, partialQtyUsed, partiallyFilled.ID)
		if err != nil {
			return err
		}
	}

	// Final update for incoming order
	if order.Type == models.MarketOrder && order.RemainingQty > 0 {
		_, err = tx.Exec(`UPDATE orders SET remaining_quantity=?, status='canceled' WHERE id=?`, order.RemainingQty, order.ID)
		order.Status = "canceled"
	} else if order.RemainingQty == 0 {
		_, err = tx.Exec(`UPDATE orders SET status='filled' WHERE id=?`, order.ID)
		order.Status = "filled"
	} else {
		_, err = tx.Exec(`UPDATE orders SET remaining_quantity=?, status='open' WHERE id=?`, order.RemainingQty, order.ID)
		order.Status = "open"
	}

	return err
}

func chooseBuyID(side models.OrderSide, incomingID, matchID int64) int64 {
	if side == models.BuySide {
		return incomingID
	}
	return matchID
}

func chooseSellID(side models.OrderSide, incomingID, matchID int64) int64 {
	if side == models.SellSide {
		return incomingID
	}
	return matchID
}

func cancelOrder(id int64) error {
	res, err := db.DB.Exec(`UPDATE orders SET status='canceled' WHERE id=? AND status='open'`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("order not found or not open")
	}
	return nil
}

func getOrderBook(symbol string) (map[string][]models.Order, error) {
	book := make(map[string][]models.Order)
	for _, side := range []models.OrderSide{"buy", "sell"} {
		order := "price DESC"
		if side == models.SellSide {
			order = "price ASC"
		}
		rows, err := db.DB.Query(`SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at
			FROM orders WHERE symbol=? AND side=? AND status='open'
			ORDER BY `+order+`, created_at ASC`, symbol, side)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var orders []models.Order
		for rows.Next() {
			var o models.Order
			rows.Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQty, &o.Status, &o.CreatedAt)
			orders = append(orders, o)
		}
		book[string(side)] = orders
	}
	return book, nil
}

func getTrades(symbol string) ([]models.Trade, error) {
	rows, err := db.DB.Query(`SELECT id, buy_order_id, sell_order_id, symbol, price, quantity, created_at FROM trades WHERE symbol=?`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trades []models.Trade
	for rows.Next() {
		var t models.Trade
		rows.Scan(&t.ID, &t.BuyOrderID, &t.SellOrderID, &t.Symbol, &t.Price, &t.Quantity, &t.CreatedAt)
		trades = append(trades, t)
	}
	return trades, nil
}

func getOrderStatus(id int64) (*models.Order, error) {
	var o models.Order
	err := db.DB.QueryRow(`SELECT id, symbol, side, type, price, quantity, remaining_quantity, status, created_at FROM orders WHERE id=?`, id).
		Scan(&o.ID, &o.Symbol, &o.Side, &o.Type, &o.Price, &o.Quantity, &o.RemainingQty, &o.Status, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	} else if err != nil {
		return nil, err
	}
	return &o, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
