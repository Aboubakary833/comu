-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id BINARY(16) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expired_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	revoked BOOLEAN NOT NULL,

	INDEX refresh_token_token_idx (token)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
