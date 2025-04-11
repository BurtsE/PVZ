CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(255)        NOT NULL
);

CREATE TYPE cities_available AS ENUM (
    'Москва',
    'Санкт-Петербург',
    'Казань'
    );

CREATE TABLE points
(
    id                UUID PRIMARY KEY,
    registration_date TIMESTAMP        NOT NULL,
    city              cities_available NOT NULL
);

CREATE TYPE status AS ENUM (
    'in_progress',
    'close'
    );

CREATE TABLE receptions
(
    id                UUID PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL,
    status            status    NOT NULL,
    point_id          UUID      NOT NULL REFERENCES points ("id")
);

CREATE TYPE product_type AS ENUM (
    'электроника',
    'одежда',
    'обувь'
    );

CREATE TABLE products
(
    id           UUID PRIMARY KEY,
    arrival_date TIMESTAMP    NOT NULL,
    type         product_type NOT NULL,
    reception_id UUID         NOT NULL REFERENCES receptions ("id")
);