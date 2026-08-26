-- +goose Up
CREATE INDEX idx_game_news_reports_game_id_created_at
    ON game_news_reports (game_id, created_at DESC);

-- +goose Down
DROP INDEX idx_game_news_reports_game_id_created_at;