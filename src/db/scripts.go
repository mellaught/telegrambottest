package db

var CREATE_USERS_IF_NOT_EXISTS = `
create table if not exists users (
	id INT NOT NULL,
   	chat_id BIGINT NOT NULL,
   	lang VARCHAR(8),
   	PRIMARY KEY (id)
);`

var CREATE_SALES_IF_NOT_EXISTS = `
create table if not exists sales (
	id SERIAL,
	user_id INT NOT NULL,
	tag VARCHAR(255),
	coin VARCHAR(255),
	price INT,
	amount VARCHAR(255),
	minter_address VARCHAR(255),
	created_at timestamp,
	last_sell_at timestamp,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users (id)
);`
