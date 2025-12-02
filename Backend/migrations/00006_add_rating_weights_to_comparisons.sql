-- +goose Up
-- +goose StatementBegin
ALTER TABLE comparison_collections 
ADD COLUMN price_rating_weight INTEGER NOT NULL DEFAULT 20,
ADD COLUMN pros_cons_rating_weight INTEGER NOT NULL DEFAULT 80;

CREATE INDEX idx_comparison_collections_rating_weights 
ON comparison_collections(price_rating_weight, pros_cons_rating_weight);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_comparison_collections_rating_weights;
ALTER TABLE comparison_collections 
DROP COLUMN IF EXISTS price_rating_weight,
DROP COLUMN IF EXISTS pros_cons_rating_weight;
-- +goose StatementEnd