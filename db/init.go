// db/init.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 3)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %v", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	DB = db
	return db, nil
}

func createTables(db *sql.DB) error {
	orderTable := `CREATE TABLE IF NOT EXISTS orders (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        symbol VARCHAR(20) NOT NULL,
        side ENUM('buy', 'sell') NOT NULL,
        type ENUM('limit', 'market') NOT NULL,
        price DECIMAL(18,8),
        quantity INT NOT NULL,
        remaining_quantity INT NOT NULL,
        status ENUM('open', 'filled', 'canceled') NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );`

	tradeTable := `CREATE TABLE IF NOT EXISTS trades (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        symbol VARCHAR(20) NOT NULL,
        buy_order_id BIGINT NOT NULL,
        sell_order_id BIGINT NOT NULL,
        price DECIMAL(18,8) NOT NULL,
        quantity INT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (buy_order_id) REFERENCES orders(id),
        FOREIGN KEY (sell_order_id) REFERENCES orders(id)
    );`

	indexQuery := `CREATE INDEX idx_symbol_side_price_time
	           ON orders(symbol, side, price, created_at)`

	if _, err := db.Exec(orderTable); err != nil {
		return fmt.Errorf("creating orders table: %w", err)
	}
	if _, err := db.Exec(tradeTable); err != nil {
		return fmt.Errorf("creating trades table: %w", err)
	}

	if _, err := db.Exec(indexQuery); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1061 {
			fmt.Println("index already exists")
			return nil
		} else {
			return fmt.Errorf("creating index: %w", err)
		}

	}
	return nil
}
