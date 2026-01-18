-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reset_tokens (
    user_id BINARY(16) NOT NULL,
	user_email VARCHAR(250) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expired_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

	INDEX reset_token_token_idx (token)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reset_tokens;
-- +goose StatementEnd
