
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `client` (
    name VARCHAR(64)
    );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `client`;
