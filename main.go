package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GOLANG-ORDER-MATCHING-SYSTEM/db"
	"github.com/GOLANG-ORDER-MATCHING-SYSTEM/router"
)

func main() {
	// Connect to MySQL
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	r := router.InitRouter()

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
