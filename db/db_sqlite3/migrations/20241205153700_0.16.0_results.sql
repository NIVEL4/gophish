
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "results"
ADD COLUMN phone_number varchar(32);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "results"
DROP COLUMN phone_number;
