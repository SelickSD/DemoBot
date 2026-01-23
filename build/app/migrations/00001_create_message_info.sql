-- +goose Up
CREATE TABLE message_info (
                              id SERIAL PRIMARY KEY,
                              message_id BIGINT NOT NULL,
                              chat_id BIGINT NOT NULL,
                              message TEXT NOT NULL,
                              user_id BIGINT NOT NULL,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE message_info;
