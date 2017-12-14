CREATE DATABASE IF NOT EXISTS workshop;

DROP TABLE IF EXISTS workshop.Workshops;

CREATE TABLE workshop.Workshops (
	id INT NOT NULL AUTO_INCREMENT,
	workshop_id VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	description VARCHAR(1024) NOT NULL,
	start_time DATETIME NOT NULL,
	end_time DATETIME NOT NULL,
	created_at DATETIME,
	updated_at DATETIME,
	cap INT UNSIGNED NOT NULL,
	cost DECIMAL(10, 2) NOT NULL,
	location VARCHAR(255) NOT NULL,
	level VARCHAR(255) NOT NULL,
	PRIMARY KEY(id),
	INDEX(name),
	INDEX(start_time),
	CONSTRAINT name_workshop_id UNIQUE(workshop_id, name)
) engine=InnoDB;
