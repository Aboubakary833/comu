-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
    id              BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    post_id         BINARY(16) NOT NULL,
    user_id         BINARY(16) NOT NULL,
    content         TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comments;
-- +goose StatementEnd
