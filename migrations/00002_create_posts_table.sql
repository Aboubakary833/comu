-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
    id              VARCHAR(255) NOT NULL PRIMARY KEY,
    user_id         VARCHAR(255) NOT NULL,
    title           VARCHAR(150) NOT NULL UNIQUE,
    slug            VARCHAR(200) NOT NULL UNIQUE,
    content         TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX post_slug_idx (slug)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
