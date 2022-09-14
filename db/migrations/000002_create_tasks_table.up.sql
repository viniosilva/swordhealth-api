CREATE TABLE tasks (
	id			int				NOT NULL	AUTO_INCREMENT,
	created_at	timestamp		NOT NULL,
	updated_at	timestamp		NOT NULL,
	deleted_at	timestamp,
    user_id     int             NOT NULL,
	summary		varchar(2500)	NOT NULL,
	status	    varchar(20)	    NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
