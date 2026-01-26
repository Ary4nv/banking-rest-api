package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const Addr = ":3000"

type Account struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type input struct {
	Name string `json:"name"`
}

type DepositInput struct {
	Amount int `json:"amount"`
}

type WithdrawInput struct {
	Amount int `json:"amount"`
}

type Transfer struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Amount int `json:"amount"`
}

// In-memory storage for now (will be replaced by DB later)
var accounts = map[int]Account{}

// openDB initializes and verifies a PostgreSQL connection pool
// *sql.DB is NOT a single connection — it is a pool manager
func openDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Pool settings (safe defaults)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Verify connection now (fail fast)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// helper: read {id} from URL and convert to int
func parseID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

func writeJSONError(w http.ResponseWriter, status int, message string){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Error" : message,
	})
}

func main() {
	// DB connection check (server won't start if DB is down)
	db, err := openDB()
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// Router and middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Health endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Home
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "Home Page",
		})
	})

	// GET /accounts (list all) — NOW READS FROM POSTGRES
	router.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT id, name,balance FROM accounts ORDER BY id")
		if err != nil {
			http.Error(w, "Query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		list := []Account{}
		for rows.Next() {
			var acc Account
			err := rows.Scan(&acc.ID, &acc.Name, &acc.Balance)
			if err != nil {
				http.Error(w, "database scan error", http.StatusInternalServerError)
				return
			}
			list = append(list, acc)

		}
		if err := rows.Err(); err != nil {
			http.Error(w, "database row error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string][]Account{
			"accounts": list,
		})

	})

	// GET /accounts/{id} — still in-memory for now (Step 2 will move it to DB)
	router.Get("/accounts/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//Conver the id client enters (string) to int and return error
		id, err := parseID(r)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json")
			return
		}

		//Get a row from databse and place it in acc variable
		var acc Account
		row := db.QueryRow("SELECT id, name, balance FROM accounts WHERE id = $1", id)
		err = row.Scan(&acc.ID, &acc.Name, &acc.Balance)
		//Scan of database happend but now row came back
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "account not found")	
			return

		}
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "database error")
			return

		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acc)

	})

	// POST /accounts (create) 
	router.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var in input
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil{
			writeJSONError(w, http.StatusBadRequest, "invalid json")
			return
		}

		if in.Name == ""{
			writeJSONError(w, http.StatusBadRequest, "name required")
			return
		}
		var acc Account
		err = db.QueryRow("INSERT INTO accounts (name, balance) VALUES ($1, 0) RETURNING id, name, balance;",in.Name).Scan(&acc.ID, &acc.Name, &acc.Balance)
		
		if err != nil{
			writeJSONError(w, http.StatusInternalServerError, "database error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(acc)
	})

	// POST /accounts/{id}/deposit — still in-memory for now
	router.Post("/accounts/{id}/deposit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := parseID(r)
		if err != nil{
			writeJSONError(w, http.StatusBadRequest, "invalid id")
			return
		}

		var dep DepositInput
		err = json.NewDecoder(r.Body).Decode(&dep)
		if err != nil{
			writeJSONError(w, http.StatusBadRequest, "invlaid JSON")
			return
		}

		if dep.Amount <= 0 {
			writeJSONError(w, http.StatusBadRequest, "amount must be higher than 0")
			return
		}

		var acc Account
		err = db.QueryRow("UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING id, name, balance;", dep.Amount, id ).Scan(&acc.ID, &acc.Name, &acc.Balance)
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "account not found")
			return

		}

		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "database error")
			return
		}
		
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acc)


	})

	// POST /accounts/{id}/withdraw — still in-memory for now
	router.Post("/accounts/{id}/withdraw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := parseID(r)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		acc, exist := accounts[id]
		if !exist {
			http.Error(w, "account not found", http.StatusNotFound)
			return
		}

		var wDraw WithdrawInput
		if err := json.NewDecoder(r.Body).Decode(&wDraw); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if wDraw.Amount <= 0 {
			http.Error(w, "amount must be > 0", http.StatusBadRequest)
			return
		}
		if wDraw.Amount > acc.Balance {
			http.Error(w, "insufficient funds", http.StatusBadRequest)
			return
		}

		acc.Balance -= wDraw.Amount
		accounts[id] = acc

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acc)
	})

	// POST /transfer — still in-memory for now
	router.Post("/transfer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var t Transfer
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if t.Amount <= 0 {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}
		if t.From <= 0 || t.To <= 0 {
			http.Error(w, "account id must be positive", http.StatusBadRequest)
			return
		}
		if t.From == t.To {
			http.Error(w, "cant transfer between same accounts", http.StatusBadRequest)
			return
		}

		accFrom, ok := accounts[t.From]
		if !ok {
			http.Error(w, "from account not found", http.StatusNotFound)
			return
		}
		accTo, ok := accounts[t.To]
		if !ok {
			http.Error(w, "destination account not found", http.StatusNotFound)
			return
		}

		if accFrom.Balance < t.Amount {
			http.Error(w, "insufficient balance", http.StatusBadRequest)
			return
		}

		accFrom.Balance -= t.Amount
		accTo.Balance += t.Amount

		accounts[t.From] = accFrom
		accounts[t.To] = accTo

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]Account{
			"from": accFrom,
			"to":   accTo,
		})
	})

	// Start server
	log.Println("Starting Server on ", Addr)
	log.Fatal(http.ListenAndServe(Addr, router))
}
