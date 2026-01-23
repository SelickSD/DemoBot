package messageinfo

import "time"

type MessageInfo struct {
	ID        int64
	MessageID int64
	ChatID    int64
	Message   string
	UserID    int64
	CreatedAt time.Time
}
