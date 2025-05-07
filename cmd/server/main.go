package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.com/digineat/go-broker-test/internal/model"
	"gitlab.com/digineat/go-broker-test/internal/repository"
	"gitlab.com/digineat/go-broker-test/internal/validator"
	"gitlab.com/digineat/go-broker-test/pkg/database"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	listenAddr := flag.String("listen", "8080", "HTTP server listen address")
	flag.Parse()

	// Initialize database connection
	db, err := database.InitDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize HTTP server
	mux := http.NewServeMux()
	tradeRepo := repository.NewTradeRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// POST /trades endpoint
	mux.HandleFunc("POST /trades", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Print("Invalid method ", "excepted: ", http.MethodPost, " have: ", r.Method)
			w.WriteHeader(http.StatusBadRequest)
		}

		if r.Header.Get("content-type") != "application/json" {
			log.Print("Invalid content type", " excepted: ", "application/json", " have: ", r.Header.Get("content-type"))
			w.WriteHeader(http.StatusBadRequest)
		}

		var req model.TradeRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Print("Failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
		}

		if err = validator.ValidateTrade(req.Account, req.Symbol, req.Side, req.Volume, req.Open, req.Close); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
		}

		trade := model.Trade{
			Account: req.Account,
			Symbol:  req.Symbol,
			Volume:  req.Volume,
			Open:    req.Open,
			Close:   req.Close,
			Side:    req.Side,
		}

		if err := tradeRepo.InsertTrade(r.Context(), &trade); err != nil {
			log.Printf("Error inserting trade: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		log.Print("Trade successfully inserted")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Trade queued successfully"))
	})

	// GET /stats/{acc} endpoint
	mux.HandleFunc("GET /stats/{acc}", func(w http.ResponseWriter, r *http.Request) {

		account := r.PathValue("acc")
		if account == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stats, err := statsRepo.GetAccountStats(r.Context(), account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// GET /healthz endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Data Base connection failed"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", *listenAddr)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
