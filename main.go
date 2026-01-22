package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const Addr = ":3000"

type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type input struct {
	Name string `json:"name"`
}

type DepositInput struct{
	Amount int `json:"amount"`
}

type WithdrawInput struct{
	Amount int `json:"amount"`
}

var accounts = map[string]Account{}
var nextID = 1



func main() {
	//Router and middleware Created
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	//Endpoints created
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Home Page",
		})
	})

	router.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		list := []Account{}
		for _,acc := range accounts{
			list = append(list,acc)
		}

		_ = json.NewEncoder(w).Encode(map[string][]Account{
			"accounts" : list,
		})


	})

	router.Get("/accounts/{id}",func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		
		id := chi.URLParam(r,"id")
		acc, found := accounts[id]
		if !found {
			http.Error(w,"user not found",http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acc)

		
	})

	router.Post("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var in input

		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if in.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		id := fmt.Sprintf("a%d", nextID)
		nextID++
		acc := Account{
			ID:      id,
			Name:    in.Name,
			Balance: 0,
		}

		accounts[id] = acc

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(acc)

	})

	router.Post("/accounts/{id}/deposit",func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")

		id := chi.URLParam(r,"id")
		acc, exist := accounts[id]
		if !exist{
			http.Error(w,"account not found",http.StatusNotFound)
			return
		}
		
		var dep DepositInput
		err := json.NewDecoder(r.Body).Decode(&dep)
		if err !=nil {
			http.Error(w,"invalid JSON",http.StatusBadRequest)
			return
		}
		if dep.Amount <= 0 {
			http.Error(w,"amount must be >0",http.StatusBadRequest)
			return

		}

		acc.Balance += dep.Amount
		accounts[id] = acc

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(acc)


	})

	router.Post("/accounts/{id}/withdraw",func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")

		id := chi.URLParam(r,"id")
		acc, exist:= accounts[id]
		if !exist {
			http.Error(w,"account Not Found",http.StatusNotFound)
			return
		}

		var wDraw WithdrawInput
		err := json.NewDecoder(r.Body).Decode(&wDraw)
		if err !=nil {
			http.Error(w,"invalid Json",http.StatusBadRequest)
			return
		}

		if wDraw.Amount <=0 {
			http.Error(w,"amount invalid : cant be 0 or lower",http.StatusBadRequest)
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

	//Server created and run
	log.Println("Starting Server on ", Addr)
	log.Fatal(http.ListenAndServe(Addr, router))

}
