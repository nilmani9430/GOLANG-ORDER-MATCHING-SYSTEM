# 🛠️ Golang Order Matching System

A simple order matching engine built with **Go** that simulates the core functionality of a stock exchange. This backend service processes buy/sell orders, matches them using price-time priority, records trades, and stores all data in a MySQL-compatible database using raw SQL (no ORM).

---

## 🚀 Features

- Supports **Limit** and **Market** orders  
- Matches based on **price-time priority**  
- Partially filled orders are handled correctly  
- In-memory order book synced with MySQL  
- RESTful API with clear endpoints  
- Error handling and input validation  
- Persistent logging of orders and trades  

---

## 🧰 Tech Stack

- **Language**: Go (latest stable version)  
- **Database**: MySQL
- **Framework**: HTTP framework 
- **SQL**: Raw queries (no ORM)  
- **Environment Management**: `.env` file (use [godotenv](https://github.com/joho/godotenv))  

---

## API Documentation

For detailed information about the API endpoints and how to use them, refer to the Postman documentation at [Postman Documentation](https://documenter.getpostman.com/view/30464667/2sB2qfBKLt).

---

## ⚙️ Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/GOLANG-ORDER-MATCHING-SYSTEM.git
cd GOLANG-ORDER-MATCHING-SYSTEM
```


### 2.  Environment Variables


Make sure to set the following environment variables in your `.env` file:

```env
DB_HOST = ""
DB_PORT = ""
DB_USER = ""
DB_PASSWORD = ""
DB_NAME = ""
```

### 3.Install Dependencies

```bash
go mod tidy
```


### 4. Run the server
```bash
go run main.go
```
