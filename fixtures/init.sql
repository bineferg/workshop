CREATE DATABASE IF NOT EXISTS workshop;

USE workshop;

DROP TABLE IF EXISTS workshop.signups;
DROP TABLE IF EXISTS workshop.workshops;
DROP TABLE IF EXISTS workshop.events;

CREATE TABLE workshop.workshops (
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


CREATE TABLE workshop.events (
	id INT NOT NULL AUTO_INCREMENT,
	event_id VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	description VARCHAR(1024) NOT NULL,
	start_time DATETIME NOT NULL,
	created_at DATETIME,
	updated_at DATETIME,
	cost DECIMAL(10, 2) NOT NULL,
	location VARCHAR(255) NOT NULL,
	PRIMARY KEY(id),
	INDEX(event_id),
	INDEX(name),
	INDEX(start_time),
	CONSTRAINT name_event_id UNIQUE(event_id, name)
) engine=InnoDB;

CREATE TABLE workshop.signups (
	id INT NOT NULL AUTO_INCREMENT,
	workshop_id VARCHAR(255) NOT NULL,
	first_name VARCHAR(255) NOT NULL,
	last_name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	created_at DATETIME,
	updated_at DATETIME,
	PRIMARY KEY(id),
	FOREIGN KEY(workshop_id) REFERENCES workshops(workshop_id) ON CASCADE DELETE,
	INDEX(workshop_id),
	CONSTRAINT name_email_workshop UNIQUE(first_name, workshop_id, email)
) engine=InnoDB;
