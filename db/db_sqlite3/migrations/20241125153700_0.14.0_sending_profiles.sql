
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE "whatsapp";

ALTER TABLE "smtp" 
ADD COLUMN number_id varchar(64);

ALTER TABLE "smtp"
ADD COLUMN auth_token varchar(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

