package repository

import (
	"context"
	"database/sql"
	"gitlab.com/digineat/go-broker-test/internal/model"
	"log"
)

type TradeRepository struct {
	Db *sql.DB
}

func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{Db: db}
}

func (h *TradeRepository) InsertTrade(ctx context.Context, trade *model.Trade) error {

	query := `
		INSERT INTO trades_q 
		(account, symbol, volume, open, close, side) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	if _, err := h.Db.ExecContext(
		ctx,
		query,
		trade.Account,
		trade.Symbol,
		trade.Volume,
		trade.Open,
		trade.Close,
		trade.Side,
	); err != nil {
		log.Printf("Error inserting trade: %v", err)
		return err
	}

	return nil

}

func (h *TradeRepository) GetPendingTrades(ctx context.Context) ([]model.Trade, error) {

	query := `
		SELECT id, account, symbol, volume, open, close, side, created_at
		FROM trades_q
		WHERE processed = FALSE
		ORDER BY created_at ASC
	`

	rows, err := h.Db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error getting pending trades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var trades []model.Trade

	for rows.Next() {
		trade := model.Trade{}
		if err = rows.Scan(
			&trade.Id,
			&trade.Account,
			&trade.Symbol,
			&trade.Volume,
			&trade.Open,
			&trade.Close,
			&trade.Side,
			&trade.CreatedAt,
		); err != nil {
			log.Printf("Error getting pending trades: %v", err)
			return nil, err
		}
		trades = append(trades, trade)

	}

	if err = rows.Err(); err != nil {
		log.Printf("Error getting pending trades: %v", err)
		return nil, err
	}

	return trades, nil

}

func (h *TradeRepository) MarkTradeAsProcessed(ctx context.Context, id int64) error {
	query := `UPDATE trades_q SET processed = TRUE WHERE id = ?`
	if _, err := h.Db.ExecContext(ctx, query, id); err != nil {
		log.Printf("Error marking trade as processed: %v", err)
		return err
	}

	return nil
}
