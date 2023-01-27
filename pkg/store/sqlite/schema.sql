--
-- Schema for sqlite store
--

CREATE TABLE IF NOT EXISTS items (
	id          TEXT NOT NULL PRIMARY KEY,
	name        TEXT NOT NULL,
	position    INTEGER NOT NULL,
	tag         TEXT NOT NULL,
	status      BOOLEAN NOT NULL
);