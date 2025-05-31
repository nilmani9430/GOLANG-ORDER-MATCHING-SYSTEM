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

# Testing Instruction (Choose any one from below)
### 1. Local Setup -> Follow the Setup Instruction.
### 2. Directly Pull and run the Docker Image -> Follow Instructions for running the Docker Image

---

## ⚙️ Setup Instructions
This method requires GO setup on your local computer

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

### 5. Use the Postman API documentation to test the endpoints
Link to API documentation is Present above in API Documentation section.

---

## ⚙️ Instructions for running the Docker Image

### 1.  Environment Variables
🔖 Note:
If you don’t want to use your own .env file, the Docker image will use inbuilt environment variables for testing purposes.
However, if you prefer to use your custom .env file, you can do so by mounting it when running the container<br><br>
Make sure to set the following environment variables in your `.env` file:

```env
DB_HOST = ""
DB_PORT = ""
DB_USER = ""
DB_PASSWORD = ""
DB_NAME = ""
```
### 2. Spin the container with image

#### Image name - codeblooded9/golang-order-matching-system

Run below command if you want to use your custom `.env` file 
```bash
docker run -it --env-file=.env -p 8080:8080 codeblooded9/golang-order-matching-system
```

Run below command if you want to go with predefined enviroment variables for testing purpose
```bash
docker run -it -p 8080:8080 codeblooded9/golang-order-matching-system
```
### 3. Use the Postman API documentation to test the endpoints
Link to API documentation is Present above in API Documentation section.











