-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "client_history" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" VARCHAR(64) NOT NULL DEFAULT 'Not provided',
    "email" VARCHAR(255) NOT NULL DEFAULT 'Not provided',
    "monitor_url" VARCHAR(255) NOT NULL DEFAULT 'Not provided',
    "monitor_password" VARCHAR(255) NOT NULL DEFAULT 'Not provided',
    "apolo_api_key" VARCHAR(255) NOT NULL DEFAULT 'Not provided',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "send_date" TIMESTAMP NULL,
    "change_date" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "client_history";
