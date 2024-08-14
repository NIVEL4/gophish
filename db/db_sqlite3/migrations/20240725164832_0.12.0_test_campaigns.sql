
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE campaigns ADD COLUMN is_test boolean DEFAULT false;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

