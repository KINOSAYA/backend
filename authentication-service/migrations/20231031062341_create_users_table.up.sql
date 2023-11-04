-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
                       id serial PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMPTZ,
                       updated_at TIMESTAMPTZ
);