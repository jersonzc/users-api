CREATE TABLE users
(
    id         CHARACTER(36) PRIMARY KEY,
    name       TEXT    NOT NULL,
    birth      DATE,
    email      TEXT UNIQUE,
    location   TEXT,
    created_at DATE    NOT NULL,
    updated_at DATE    NOT NULL,
    active     BOOLEAN NOT NULL
);