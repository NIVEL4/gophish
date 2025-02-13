
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "beef" (
    "url" varchar(256)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "beef";
