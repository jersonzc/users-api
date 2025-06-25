CREATE TABLE users
(
    id       INTEGER UNIQUE NOT NULL,
    name     VARCHAR(200)   NOT NULL,
    birth    DATE           NOT NULL,
    active   BOOLEAN        NOT NULL,
    location VARCHAR(200)
);