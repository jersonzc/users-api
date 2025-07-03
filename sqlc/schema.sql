CREATE TABLE users (
    id         CHARACTER(36) PRIMARY KEY,
    name       TEXT      NOT NULL,
    birth      DATE,
    email      TEXT UNIQUE,
    location   TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    active     BOOLEAN   NOT NULL
);
