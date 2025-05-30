package router

import (
	"github.com/Ringover_assignment/handlers"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Methods("GET", "POST", "PUT", "DELETE", "PATCH").Subrouter()
	subRouter.HandleFunc("/orders", handlers.HandlePlaceOrders).Methods("POST")
	subRouter.HandleFunc("/orders/{id}", handlers.HandleCancelOrder).Methods("DELETE")
	subRouter.HandleFunc("/orderbook", handlers.HandleGetOrderBook).Methods("GET")
	subRouter.HandleFunc("/trades", handlers.HandleGetTrades).Methods("GET")
	subRouter.HandleFunc("/orders/{id}", handlers.HandleGetOrderStatus).Methods("GET")

	return router
}
