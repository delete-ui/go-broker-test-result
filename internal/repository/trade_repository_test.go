package repository

import (
	"context"
	"gitlab.com/digineat/go-broker-test/internal/model"
	database2 "gitlab.com/digineat/go-broker-test/pkg/database"
	"os"
	"testing"
)

func TestTradeRepository(t *testing.T) {

	tmpDB, err := os.CreateTemp("", "testdb*-.db")
	if err != nil {
		t.Fatal("Error creating temp DB:", err)
	}
	defer os.Remove(tmpDB.Name())

	database, err := database2.InitDB(tmpDB.Name())
	if err != nil {
		t.Fatal("Error initializing DB:", err)
	}
	defer database.Close()

	repo := NewTradeRepository(database)
	ctx := context.Background()

	trade := &model.Trade{
		Account: "test123",
		Symbol:  "EURUSD",
		Volume:  1.0,
		Open:    1.1,
		Close:   1.2,
		Side:    "buy",
	}

	err = repo.InsertTrade(ctx, trade)
	if err != nil {
		t.Fatal("Error inserting trade: ", err)
	}

	trades, err := repo.GetPendingTrades(ctx)
	if err != nil {
		t.Fatal("Error getting pending trades: ", err)
	}
	if len(trades) != 1 {
		t.Fatal("Expected 1 trade, got ", len(trades))
	}

	err = repo.MarkTradeAsProcessed(ctx, trade.Id)
	if err != nil {
		t.Fatal("Error marking trade: ", err)
	}

}
