-- +goose Up
-- +goose StatementBegin
ALTER TABLE collections ADD COLUMN public_link VARCHAR(255) UNIQUE NULL;
CREATE INDEX idx_collections_public_link ON collections(public_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_collections_public_link;
ALTER TABLE collections DROP COLUMN IF EXISTS public_link;
-- +goose StatementEnd