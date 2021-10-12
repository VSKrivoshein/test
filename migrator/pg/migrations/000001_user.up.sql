CREATE TABLE users
(
    id            uuid PRIMARY KEY NOT NULL UNIQUE,
    email         varchar(255)     NOT NULL UNIQUE,
    password_hash varchar(255)     NOT NULL UNIQUE
);