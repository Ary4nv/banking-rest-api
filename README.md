# 🏦 Banking REST API (Go)

A simple **RESTful banking API** built with **Go** using `net/http` and the **Chi router**.  
This project simulates core banking operations such as account creation, deposits, withdrawals, and transfers.

The purpose of this project is to practice and demonstrate **backend fundamentals**, clean API design, and server-side logic in Go.

---

## 🚀 Features

- Create bank accounts
- List all accounts
- Get account by ID
- Deposit money
- Withdraw money (with balance validation)
- Transfer money between accounts
- JSON request / response handling
- Proper HTTP status codes
- Middleware logging & recovery

---

## 🛠 Tech Stack

- **Go**
- **net/http**
- **Chi Router**
- **JSON**
- In-memory data storage (map)

---

## 📦 Getting Started

### Prerequisites
- Go 1.20 or newer

### Run the server
```bash
go run main.go
```

The server will start on:
```
http://localhost:3000
```

---

## 🔍 API Endpoints

### Health Check
**GET** `/health`

```bash
curl http://localhost:3000/health
```

---

### Create Account
**POST** `/accounts`

```bash
curl -X POST http://localhost:3000/accounts \
  -H "Content-Type: application/json" \
  -d '{"name":"Aryan"}'
```

Response:
```json
{
  "id": "a1",
  "name": "Aryan",
  "balance": 0
}
```

---

### List All Accounts
**GET** `/accounts`

```bash
curl http://localhost:3000/accounts
```

Response:
```json
{
  "accounts": [
    {
      "id": "a1",
      "name": "Aryan",
      "balance": 100
    }
  ]
}
```

---

### Get Account by ID
**GET** `/accounts/{id}`

```bash
curl http://localhost:3000/accounts/a1
```

---

### Deposit Money
**POST** `/accounts/{id}/deposit`

```bash
curl -X POST http://localhost:3000/accounts/a1/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount":100}'
```

---

### Withdraw Money
**POST** `/accounts/{id}/withdraw`

```bash
curl -X POST http://localhost:3000/accounts/a1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount":50}'
```

Validation:
- Amount must be greater than 0
- Balance must be sufficient

---

### Transfer Money
**POST** `/transfer`

```bash
curl -X POST http://localhost:3000/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"a1","to":"a2","amount":25}'
```

Response:
```json
{
  "from": {
    "id": "a1",
    "name": "Aryan",
    "balance": 75
  },
  "to": {
    "id": "a2",
    "name": "Nima",
    "balance": 25
  }
}
```

Validation:
- Accounts must exist
- Amount must be greater than 0
- Cannot transfer to the same account
- Sender must have sufficient balance

---

## ⚠️ Important Notes

- Data is stored **in memory** using a Go map  
- Restarting the server will reset all accounts
- This is intentional for learning purposes

---

## 🧭 Roadmap / Next Steps

- Replace in-memory storage with **SQLite**
- Add unit tests
- Improve error response structure
- Add authentication
- Deploy to cloud (Render / AWS)

---

## 👤 Author

Built by **Aryan**  
Learning backend development with Go and REST APIs.

---

## 📄 License

This project is open-source and free to use for learning and practice.