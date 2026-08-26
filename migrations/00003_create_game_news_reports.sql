-- +goose Up
CREATE TABLE game_news_reports (
                                   id         BIGSERIAL PRIMARY KEY,
                                   game_id    BIGINT NOT NULL,
                                   report     TEXT NOT NULL,
                                   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE game_news_reports;
