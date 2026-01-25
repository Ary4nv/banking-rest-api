# 🏦 Banking REST API (Go)

A clean, RESTful banking API built in **Go** using the **Chi router**.  
Simulates core banking operations: account creation, balance retrieval, deposits, withdrawals, and transfers — with input validation and proper HTTP responses.

**Goal**: Demonstrate backend fundamentals, clean API design, error handling, and real-world server logic — aligned with internship-level expectations at banks, fintech, and tech companies.

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![Chi](https://img.shields.io/badge/Router-Chi-00ADD8?style=flat)](https://github.com/go-chi/chi)

---

## 🚀 Features

- Create bank accounts
- List all accounts
- Get account by ID
- Deposit money (with amount validation)
- Withdraw money (with insufficient funds check)
- Transfer money between accounts (with balance & same-account validation)
- JSON request/response with proper status codes (200, 201, 400, 404)
- Middleware: logging + panic recovery

**Current Status**: In-memory storage (map) — **PostgreSQL integration in progress** (persistent storage + transactions).

---

## 🛠 Tech Stack

- **Go** (1.20+)
- **Chi router** (lightweight routing)
- **net/http** (standard library)
- JSON encoding/decoding
- In-memory storage (`map[string]Account`)
- Middleware (`middleware.Logger`, `middleware.Recoverer`)

---

## 📦 Getting Started

### Prerequisites

- Go 1.20 or higher

### Run Locally

```bash
# Clone repo
git clone https://github.com/yourusername/banking-rest-api.git
cd banking-rest-api

# Run
go run main.go

Server starts at:  
**http://localhost:3000**

---

## 🔍 API Endpoints

Base URL: `http://localhost:3000`

| Method | Endpoint                        | Description                          | Example curl Command                                                                 | Expected Response (200/201)                          |
|--------|---------------------------------|--------------------------------------|--------------------------------------------------------------------------------------|------------------------------------------------------|
| GET    | /health                         | Health check                         | `curl http://localhost:3000/health`                                                  | `{"status": "ok"}`                                   |
| POST   | /accounts                       | Create account                       | `curl -X POST http://localhost:3000/accounts -H "Content-Type: application/json" -d '{"name":"Arian"}'` | `{"id":"a1","name":"Arian","balance":0}`             |
| GET    | /accounts                       | List all accounts                    | `curl http://localhost:3000/accounts`                                                | `[{"id":"a1","name":"Arian","balance":0}, ...]`      |
| GET    | /accounts/{id}                  | Get account by ID                    | `curl http://localhost:3000/accounts/a1`                                             | `{"id":"a1","name":"Arian","balance":0}`             |
| POST   | /accounts/{id}/deposit          | Deposit money                        | `curl -X POST http://localhost:3000/accounts/a1/deposit -H "Content-Type: application/json" -d '{"amount":100}'` | `{"id":"a1","name":"Arian","balance":100}`           |
| POST   | /accounts/{id}/withdraw         | Withdraw money                       | `curl -X POST http://localhost:3000/accounts/a1/withdraw -H "Content-Type: application/json" -d '{"amount":50}'` | `{"id":"a1","name":"Arian","balance":50}`            |
| POST   | /transfer                       | Transfer money between accounts      | `curl -X POST http://localhost:3000/transfer -H "Content-Type: application/json" -d '{"from":"a1","to":"a2","amount":25}'` | `{"from":{...},"to":{...}}`                          |

**Error Responses** (examples):  
- 400 Bad Request: `{"error": "amount must be > 0"}`  
- 404 Not Found: `{"error": "account not found"}`  
- 400 Bad Request: `{"error": "insufficient funds"}`

---

## ⚠️ Current Limitations

- Data is **in-memory only** (resets on restart) — PostgreSQL integration in progress
- No authentication — learning/demo purposes only
- No unit tests yet — planned next

---

## 🧭 Roadmap / Next Steps

- Integrate PostgreSQL for persistent storage
- Add transactional safety for transfers
- Containerize with Docker
- Add minimal frontend demo (HTML + JS fetch)
- Write unit tests for handlers
- Deploy to Render (public live demo)
- Prepare for interviews (explain endpoints, tradeoffs, errors)

---

## 👤 Author

Built by **Arian Vares**  
Fourth-year Computer Science student at Ontario Tech University  
Learning backend development with Go, REST APIs, and databases.

Open-source for learning and practice — feel free to fork/use.

```

