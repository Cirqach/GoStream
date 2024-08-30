 CREATE TABLE queue (
	id		SERIAL PRIMARY KEY,
	path		VARCHAR(255) NOT NULL,
	broadcast_time	TIMESTAMP NOT NULL
);
DROP TABLE IF EXISTS queue;
 CREATE TABLE queue (
	id		SERIAL PRIMARY KEY,
	path		VARCHAR(255) NOT NULL,
	broadcast_time	TIMESTAMP NOT NULL
);
INSERT INTO queue
	(path, broadcast_time)
VALUES
	('./video/unprocessed/2.mp4','2024-08-07 12:10:00' ),
	('./video/unprocessed/1.mp4','2024-10-07 11:30:00' ),
	('./video/unprocessed/1.mp4','2024-12-07 11:30:00' ),
	('./video/unprocessed/1.mp4','2024-08-07 11:30:00' ),
	('./video/unprocessed/1.mp4','2024-08-07 11:30:00' ),
	('1', '2024-08-10 08:20:00'),
	('./video/unprocessed/3.mp4','2024-08-07 13:15:00' ),
	('./video/unprocessed/4.mp4','2024-08-07 13:30:00' ),
	('./video/unprocessed/5.mp4','2024-08-07 15:30:00' );

DROP TABLE IF EXISTS users;
CREATE TABLE users (
	id		SERIAL PRIMARY KEY,
	username	VARCHAR(255) NOT NULL,
	password	VARCHAR(255) NOT NULL,
	email		VARCHAR(255) NOT NULL
);

INSERT INTO users
	(username, password)
VALUES
	('admin','$2a$12$ZdsWrIAfaRjNhhOHMgX6GOiBqxnqnIgN.coGSjL3AQ2McRjg5SmLS'),
	('user','$2a$12$CsLeQ75XnprA5N53OarPpuH0MBYKonYJzkkci0jMv7eFbkYicJ4S6');

