-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS otp_codes (
    type INTEGER NOT NULL,
    user_email VARCHAR(250) NOT NULL,
    value VARCHAR(255) NOT NULL UNIQUE,
    expired_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX otp_code_value_idx (value),
    INDEX otp_code_user_email_idx (user_email)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE otp_codes;
-- +goose StatementEnd