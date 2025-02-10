-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE "client" ADD COLUMN "email" VARCHAR(255) NOT NULL DEFAULT 'Not provided';

ALTER TABLE "client" ADD COLUMN "monitor_url" VARCHAR(255) NOT NULL DEFAULT 'Not provided';

ALTER TABLE "client" ADD COLUMN "monitor_password" VARCHAR(255) NOT NULL DEFAULT 'Not provided';

ALTER TABLE "client" ADD COLUMN "apolo_api_key" VARCHAR(255) NOT NULL DEFAULT 'Not provided';

ALTER TABLE "client" ADD COLUMN "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE "client" ADD COLUMN "send_date" TIMESTAMP NULL;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "client" DROP COLUMN "email";

ALTER TABLE "client" DROP COLUMN "monitor_url";

ALTER TABLE "client" DROP COLUMN "monitor_password";

ALTER TABLE "client" DROP COLUMN "apolo_api_key";

ALTER TABLE "client" DROP COLUMN "created_at";

ALTER TABLE "client" DROP COLUMN "send_date";
