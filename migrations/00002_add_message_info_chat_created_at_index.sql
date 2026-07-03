-- +goose Up
CREATE INDEX idx_message_info_chat_created_at
    ON message_info (chat_id, created_at DESC);

-- +goose Down
DROP INDEX idx_message_info_chat_created_at;
