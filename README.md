```markdown
# 🏦 Banking REST API (Go)

A clean, RESTful banking API built in **Go** using the **Chi router** and **PostgreSQL**.  
Simulates core banking operations: account creation, balance retrieval, deposits, withdrawals, and transfers — with input validation, proper HTTP responses, and real database persistence.

**Goal**: Demonstrate backend fundamentals, clean API design, SQL safety, error handling, and real-world server logic — aligned with internship-level expectations at banks, fintech, and tech companies.

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![Chi](https://img.shields.io/badge/Router-Chi-00ADD8?style=flat)](https://github.com/go-chi/chi)
[![PostgreSQL](https://img.shields.io/badge/DB-PostgreSQL-336791?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com)

---

## 🚀 Features

- Create bank accounts (DB-backed with auto-generated IDs)
- List all accounts (DB-backed)
- Get account by ID (DB-backed with proper 404 handling)
- Deposit money (with amount validation, DB-backed)
- Withdraw money (with insufficient funds check, DB-backed)
- Transfer money between accounts (with balance & same-account validation, DB transaction)
- JSON request/response with proper status codes (200, 201, 400, 404)
- Middleware: logging + panic recovery
- PostgreSQL persistence (data survives restarts)
- Docker containerization (consistent local runs)

**Current Status**: Full PostgreSQL integration complete — all endpoints (create, list, get, deposit, withdraw, transfer with transaction) now use real DB persistence.

---

## 🛠 Tech Stack

- **Go** (1.20+)
- **Chi router** (lightweight routing)
- **net/http** (standard library)
- **PostgreSQL** (persistent storage)
- **database/sql** + **pgx** driver (parameterized queries with RETURNING)
- JSON encoding/decoding
- Middleware (`middleware.Logger`, `middleware.Recoverer`)
- **Docker** (containerization)

---

## 📦 Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL running locally (or Docker Compose)
- Docker (optional, recommended)

### Database Setup (One-Time – if not using Docker)

```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 0
);
```

### Environment Variable

Set your database connection string:

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/banking_api?sslmode=disable"
```

### Run Locally (without Docker)

```bash
# Clone repo
git clone https://github.com/Ary4nv/banking-rest-api.git
cd banking-rest-api

# Run
go run main.go
```

Server starts at:  
**http://localhost:3000**

### Run with Docker Compose (recommended)

```bash
# Start Postgres + API
docker-compose up --build
```

- App runs at **http://localhost:8080**  
- Postgres exposed at **localhost:5432** (user: appuser, password: secret, db: dbapp)

---

## 🔍 API Endpoints

Base URL: `http://localhost:3000` (or :8080 with Docker)

| Method | Endpoint                        | Description                          | Example curl Command                                                                 | Expected Response (200/201)                          |
|--------|---------------------------------|--------------------------------------|--------------------------------------------------------------------------------------|------------------------------------------------------|
| GET    | /health                         | Health check                         | `curl http://localhost:3000/health`                                                  | `{"status": "ok"}`                                   |
| POST   | /accounts                       | Create account                       | `curl -X POST http://localhost:3000/accounts -H "Content-Type: application/json" -d '{"name":"Arian"}'` | `{"id":1,"name":"Arian","balance":0}`                |
| GET    | /accounts                       | List all accounts                    | `curl http://localhost:3000/accounts`                                                | `[{"id":1,"name":"Arian","balance":0}, ...]`         |
| GET    | /accounts/{id}                  | Get account by ID                    | `curl http://localhost:3000/accounts/1`                                              | `{"id":1,"name":"Arian","balance":0}`                |
| POST   | /accounts/{id}/deposit          | Deposit money                        | `curl -X POST http://localhost:3000/accounts/1/deposit -H "Content-Type: application/json" -d '{"amount":100}'` | `{"id":1,"name":"Arian","balance":100}`              |
| POST   | /accounts/{id}/withdraw         | Withdraw money                       | `curl -X POST http://localhost:3000/accounts/1/withdraw -H "Content-Type: application/json" -d '{"amount":50}'` | `{"id":1,"name":"Arian","balance":50}`               |
| POST   | /transfer                       | Transfer money                       | `curl -X POST http://localhost:3000/transfer -H "Content-Type: application/json" -d '{"from":1,"to":2,"amount":25}'` | `{"from":{...},"to":{...}}`                          |

**Error Responses** (examples):  
- 400 Bad Request: `{"Error": "amount must be > 0"}`  
- 404 Not Found: `{"Error": "account not found"}`  
- 400 Bad Request: `{"Error": "insufficient funds"}`

---

## ⚠️ Current Limitations

- No authentication — learning/demo purposes only
- No unit tests yet — planned next

---

## 🧭 Roadmap / Next Steps

- Add unit tests for handlers
- Add minimal frontend demo (HTML + JS fetch)
- Deploy to Render (public live demo)
- Prepare for interviews (explain endpoints, tradeoffs, errors)

---

## 👤 Author

Built by **Arian Vares**  
Fourth-year Computer Science student at Ontario Tech University  
Learning backend development with Go, REST APIs, and databases.

Open-source for learning and practice — feel free to fork/use.
