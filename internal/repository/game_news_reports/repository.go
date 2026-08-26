package game_news_reports

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameNewsReportsRepository struct {
	db *pgxpool.Pool
}

func NewGameNewsReportsRepository(db *pgxpool.Pool) *GameNewsReportsRepository {
	return &GameNewsReportsRepository{
		db: db,
	}
}

func (r *GameNewsReportsRepository) CreateGameNewsReport(
	ctx context.Context,
	gameID int64,
	report string,
) error {
	const query = `
        INSERT INTO game_news_reports (game_id, report)
        VALUES ($1, $2)
    `

	_, err := r.db.Exec(ctx, query, gameID, report)
	if err != nil {
		return fmt.Errorf("create game news report: %w", err)
	}

	return nil
}

func (r *GameNewsReportsRepository) GetLatestGameNewsReport(
	ctx context.Context,
	gameID int64,
) (*GameNewsReport, error) {
	const query = `
        SELECT id, game_id, report, created_at
        FROM game_news_reports
        WHERE game_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `

	var result GameNewsReport

	err := r.db.QueryRow(ctx, query, gameID).Scan(
		&result.ID,
		&result.GameID,
		&result.Report,
		&result.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get latest game news report: %w", err)
	}

	return &result, nil
}
