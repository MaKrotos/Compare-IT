-- +goose Up
-- +goose StatementBegin
CREATE TABLE collections (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES telegram_users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_collections_user_id ON collections(user_id);
CREATE INDEX idx_collections_is_pinned ON collections(is_pinned);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE collections;
-- +goose StatementEnd