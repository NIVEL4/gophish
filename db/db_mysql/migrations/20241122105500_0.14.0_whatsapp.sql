
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `whatsapp` (
	id integer primary key auto_increment,
	user_id bigint,
	interface_type varchar(255),
	name varchar(255),
	number varchar(255),
	auth_token varchar(255),
	modified_date datetime
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `whatsapp`;
