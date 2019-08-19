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
	id BIGINT NOT NULL,
	user_id INT NOT NULL,
	info VARCHAR(255),
	sell_at timestamp,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users (id)
);`
