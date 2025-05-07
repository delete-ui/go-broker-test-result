package repository

import (
	"context"
	"database/sql"
	"gitlab.com/digineat/go-broker-test/internal/model"
	"log"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (h *StatsRepository) GetAccountStats(ctx context.Context, account string) (*model.AccountStats, error) {

	query := `
		SELECT account, trades, profit
		FROM account_stats
		WHERE account = ?
	`

	row := h.db.QueryRowContext(ctx, query, account)
	var stat model.AccountStats
	if err := row.Scan(&stat.Account, &stat.Trades, &stat.Profit); err != nil {
		log.Printf("Failed to get account stats: %v", err)
		return nil, err
	}

	return &stat, nil

}

func (h *StatsRepository) UpdateAccountStats(ctx context.Context, stats *model.AccountStats) error {

	query := `
		INSERT INTO account_stats (account, trades, profit)
		VALUES (?, ?, ?)
		ON CONFLICT(account) DO UPDATE SET
			trades = excluded.trades,
			profit = excluded.profit
	`

	if _, err := h.db.ExecContext(ctx, query, stats.Account, stats.Trades, stats.Profit); err != nil {
		log.Printf("Failed to update account stats: %v", err)
		return err
	}

	return nil

}
