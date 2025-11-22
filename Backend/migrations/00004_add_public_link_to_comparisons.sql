-- +goose Up
-- +goose StatementBegin
ALTER TABLE comparison_collections ADD COLUMN public_link VARCHAR(255) UNIQUE NULL;
CREATE INDEX idx_comparison_collections_public_link ON comparison_collections(public_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_comparison_collections_public_link;
ALTER TABLE comparison_collections DROP COLUMN IF EXISTS public_link;
-- +goose StatementEnd