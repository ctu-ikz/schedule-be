-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens
(
    id           UUID PRIMARY KEY,
    user_id      UUID                     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    hashed_token TEXT                     NOT NULL,
    expires_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked      BOOLEAN                  DEFAULT FALSE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    device_info  TEXT,
    ip_address   INET
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
