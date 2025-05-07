package main

import (
	"context"
	"flag"
	"gitlab.com/digineat/go-broker-test/internal/repository"
	"gitlab.com/digineat/go-broker-test/pkg/database"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	pollInterval := flag.Duration("poll", 100*time.Millisecond, "polling interval")
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

	tradeRepo := repository.NewTradeRepository(db)
	statsRepository := repository.NewStatsRepository(db)

	log.Printf("Worker started with polling interval: %v", *pollInterval)

	// Main worker loop
	for {
		processTrades(context.Background(), tradeRepo, statsRepository)
		time.Sleep(*pollInterval)
	}
}

func processTrades(ctx context.Context, tradeRepo *repository.TradeRepository, statsRepo *repository.StatsRepository) {

	tx, err := tradeRepo.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback()

	trades, err := tradeRepo.GetPendingTrades(ctx)
	if err != nil {
		log.Printf("Failed to get pending trades: %v", err)
		return
	}

	if len(trades) == 0 {
		return
	}

	for _, trade := range trades {

		stats, err := statsRepo.GetAccountStats(ctx, trade.Account)
		if err != nil {
			log.Printf("Failed to get account stats: %v", err)
			return
		}

		lot := 100000.0
		profit := (trade.Close - trade.Open) * trade.Volume * lot
		if trade.Side == "sell" {
			profit = -profit
		}

		stats.Trades++
		stats.Profit += profit

		if err := statsRepo.UpdateAccountStats(ctx, stats); err != nil {
			log.Printf("Failed to update account stats: %v", err)
			return
		}

		if err := tradeRepo.MarkTradeAsProcessed(ctx, trade.Id); err != nil {
			log.Printf("Failed to mark trade as processed: %v", err)
			return
		}
		log.Printf("Processed trade ID %d for account %s, profit: %.2f", trade.Id, trade.Account, profit)

	}
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}
}
