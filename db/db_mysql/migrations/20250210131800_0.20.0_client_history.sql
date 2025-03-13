-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `client_history` (
    `id` SERIAL PRIMARY KEY,
    `name` VARCHAR(64) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `monitor_url` VARCHAR(255) NOT NULL,
    `monitor_password` VARCHAR(255) NOT NULL,
    `apolo_api_key` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `sent_date` TIMESTAMP NULL DEFAULT NULL,
    `sent_by` VARCHAR(30) NOT NULL DEFAULT 'Not sent',
    `send_method` VARCHAR(20) NOT NULL DEFAULT 'None'
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `client_history`;
