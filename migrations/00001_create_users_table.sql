-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(250) NOT NULL UNIQUE,
    email_verified_at DATETIME NULL,
    avatar VARCHAR(2083) DEFAULT "",
    active BOOLEAN DEFAULT true,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX user_email_idx (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd