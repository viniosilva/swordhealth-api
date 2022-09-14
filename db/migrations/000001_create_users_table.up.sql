CREATE TABLE users (
	id			int				NOT NULL	AUTO_INCREMENT,
	created_at	timestamp		NOT NULL,
	updated_at	timestamp		NOT NULL,
	deleted_at	timestamp,
	email		varchar(250)	NOT NULL,
	username	varchar(100)	NOT NULL,
	password	varchar(128)	NOT NULL,
	role		varchar(20)		NOT NULL,
	PRIMARY KEY (id)
);
