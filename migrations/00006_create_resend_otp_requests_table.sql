-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resend_otp_requests (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    user_email VARCHAR(250) NOT NULL,
	count INTEGER NOT NULL,
    last_sent_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

	INDEX resend_otp_request_user_email_idx (user_email)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resend_otp_requests;
-- +goose StatementEnd
