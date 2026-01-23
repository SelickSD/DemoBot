package messageinfo

import (
	"context"

	"github.com/SelickSD/DemoBot.git/internal/repository/db"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Save(ctx context.Context, msg MessageInfo) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO message_info (message_id, chat_id, message, user_id)
		VALUES ($1, $2, $3, $4)
	`,
		msg.MessageID,
		msg.ChatID,
		msg.Message,
		msg.UserID,
	)

	return err
}

func (r *Repository) GetByChatID(
	ctx context.Context,
	chatID int64,
	limit int,
) ([]MessageInfo, error) {

	rows, err := db.Pool.Query(ctx, `
		SELECT id, message_id, chat_id, message, user_id, created_at
		FROM message_info
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MessageInfo

	for rows.Next() {
		var m MessageInfo
		if err := rows.Scan(
			&m.ID,
			&m.MessageID,
			&m.ChatID,
			&m.Message,
			&m.UserID,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}

func (r *Repository) DeleteAll(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM message_info`)
	return err
}
