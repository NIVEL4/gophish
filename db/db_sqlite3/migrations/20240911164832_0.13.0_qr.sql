
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "qr_conf" (
    "user_id" integer, "size" integer DEFAULT 256, "pixels" varchar(8) DEFAULT '#000000', "background" varchar(8) DEFAULT '#ffffff'
    );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "qr_conf";
