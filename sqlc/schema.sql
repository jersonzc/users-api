CREATE TABLE users (
    id         CHARACTER(36) PRIMARY KEY,
    name       TEXT      NOT NULL,
    birth      DATE,
    email      TEXT UNIQUE,
    location   TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    active     BOOLEAN   NOT NULL
);
