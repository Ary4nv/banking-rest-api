package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestValidateTransfer_Valid(t *testing.T) {
	tr := Transfer{From: 1, To: 2, Amount: 10}

	err := validateTransfer(tr)
	if err != nil {
		t.Fatalf("expected nil error %v", err)

	}
}

func TestValidateTransfer_InvalidCases(t *testing.T) {
	tests := []struct {
		name    string
		input   Transfer
		wantErr bool
	}{
		{
			name:    "From id is 0",
			input:   Transfer{From: 0, To: 2, Amount: 10},
			wantErr: true,
		},
		{
			name:    "to id cant be negetive",
			input:   Transfer{From: 1, To: -2, Amount: 5},
			wantErr: true,
		},
		{
			name:    "From and To are same",
			input:   Transfer{From: 2, To: 2, Amount: 10},
			wantErr: true,
		},
		{
			name:    "amount can be 0",
			input:   Transfer{From: 1, To: 4, Amount: 0},
			wantErr: true,
		},
		{
			name:    "amount cant be negetive",
			input:   Transfer{From: 3, To: 2, Amount: -10},
			wantErr: true,
		}, {

			name:    "valid transfer",
			input:   Transfer{From: 1, To: 2, Amount: 10},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTransfer(tc.input)

			if tc.wantErr && err == nil {
				t.Fatal("expected error but got nil")

			}
			if !tc.wantErr && err != nil {
				t.Fatal("not expecting error but got one")

			}
		})

	}
}

// /healt endpoint test
func TestHealthEndpoint(t *testing.T) {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	rq := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 but got %d", rr.Code)
	}
	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Fatalf("decoding error : %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status : ok but got %s", body["status"])
	}
}

func TestHomeEndpoint(t *testing.T) {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "home page",
		})

	})

	rq := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected to get 200 but got %d", rr.Code)
	}

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Fatalf("faild to decode : %v", err)
	}

	if body["message"] != "home page" {
		t.Fatalf("we expected message:home page but got message:%s", body["message"])
	}

}

func TestUnknownRouteReturn404(t *testing.T) {
	router := chi.NewRouter()

	rq := httptest.NewRequest("GET", "/does-not-exist", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("we expected to get 404 but got : %d", rr.Code)
	}
}

func TestTransferReturn400(t *testing.T) {
	router := chi.NewRouter()

	router.Post("/transfer", func(w http.ResponseWriter, r *http.Request) {
		var tr Transfer
		err := json.NewDecoder(r.Body).Decode(&tr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "server failed")
	})

	rq := httptest.NewRequest("POST", "/transfer", strings.NewReader("{bad json"))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected to get 400 but got %v", rr.Code)
	}

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)

	if err != nil {
		t.Fatalf("failed to decode response : %v", err)

	}
	if body["Error"] != "invalid JSON" {
		t.Fatalf("expected Error: invalid JSON but got Error : %s", body["Error"])
	}

}

func TestFromEqualsTo(t *testing.T) {
	router := chi.NewRouter()

	router.Post("/transfer", func(w http.ResponseWriter, r *http.Request) {

		var tr Transfer
		err := json.NewDecoder(r.Body).Decode(&tr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := validateTransfer(tr); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

	})

	rq := httptest.NewRequest("POST", "/transfer", strings.NewReader(`{"from":1, "to":1, "amount":10}`))
	rq.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", rr.Code)
	}

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode : %s", rr.Body.String())
	}

	if body["Error"]!= "cannot transfer to same account"{
		t.Fatalf("expected to get error: cannot transfer to same account but got %s",body["Error"])
	}


}


func TestInvalidAmountReturns400 (t *testing.T){
	router:=chi.NewRouter()

	router.Post("/transfer", func(w http.ResponseWriter, r *http.Request){
		var tr Transfer
		err := json.NewDecoder(r.Body).Decode(&tr)
		if err != nil {
			writeJSONError (w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := validateTransfer(tr); err !=nil{
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

	})

	rq := httptest.NewRequest("POST", "/transfer", strings.NewReader(`{"from":1, "to":2, "amount":0}`))
	rq.Header.Set("Content-Type","application/json")

	rr:= httptest.NewRecorder()

	router.ServeHTTP(rr, rq)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 but got : %d", rr.Code)
	}

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	if err != nil{
		t.Fatalf("failed to decode : %s",rr.Body.String())
	}
	if body["Error"] != "amount must be greater than 0"{
		t.Fatalf("expect error : amount must be grater than zero but got %s",body["Error"])

	}

}