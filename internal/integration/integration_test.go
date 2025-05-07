package integration

import (
	"context"
	"gitlab.com/digineat/go-broker-test/internal/model"
	"gitlab.com/digineat/go-broker-test/internal/repository"
	database2 "gitlab.com/digineat/go-broker-test/pkg/database"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestIntegration(t *testing.T) {
	// Setup test database
	tmpDB, err := os.CreateTemp("", "integration-*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpDB.Name())

	database, err := database2.InitDB(tmpDB.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	// Initialize repositories
	tradeRepo := repository.NewTradeRepository(database)
	statsRepo := repository.NewStatsRepository(database)

	// Test trade processing
	trade := &model.Trade{
		Account: "test123",
		Symbol:  "EURUSD",
		Volume:  1.0,
		Open:    1.1,
		Close:   1.2,
		Side:    "buy",
	}

	// Insert trade
	if err := tradeRepo.InsertTrade(context.Background(), trade); err != nil {
		t.Fatalf("InsertTrade failed: %v", err)
	}

	// Process trade (simulate worker)
	time.Sleep(150 * time.Millisecond) // Wait for worker to process

	// Check stats
	stats, err := statsRepo.GetAccountStats(context.Background(), "test123")
	if err != nil {
		t.Fatalf("GetAccountStats failed: %v", err)
	}

	if stats.Trades != 1 {
		t.Errorf("Expected 1 trade, got %d", stats.Trades)
	}

	expectedProfit := (1.2 - 1.1) * 1.0 * 100000.0
	if stats.Profit != expectedProfit {
		t.Errorf("Expected profit %.2f, got %.2f", expectedProfit, stats.Profit)
	}
}
