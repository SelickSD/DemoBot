package game_news_reports

import "time"

type GameNewsReport struct {
	ID        int64
	GameID    int64
	Report    string
	CreatedAt time.Time
}
