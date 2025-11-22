-- +goose Up
-- +goose StatementBegin
CREATE TABLE comparison_collections (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES telegram_users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    items JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comparison_collections_user_id ON comparison_collections(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE comparison_collections;
-- +goose StatementEnd