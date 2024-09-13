
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `qr_conf` (
    size bigint DEFAULT 256, pixels varchar(8) DEFAULT '#000000', background varchar(8) DEFAULT '#ffffff'
    );

INSERT INTO `qr_conf` VALUES (256, '#000000', '#ffffff');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `qr_conf`;
