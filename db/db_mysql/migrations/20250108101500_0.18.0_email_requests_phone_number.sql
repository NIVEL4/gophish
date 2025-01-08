
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `email_requests` 
ADD COLUMN phone_number varchar(32);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `email_requests`
DROP COLUMN phone_number;
