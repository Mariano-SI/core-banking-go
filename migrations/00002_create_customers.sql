-- +goose Up
CREATE TABLE customers (
    id         uuid        PRIMARY KEY,
    name       text        NOT NULL,
    cpf        char(11)    NOT NULL UNIQUE,
    email      text        NOT NULL UNIQUE,
    birth_date date,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE customers;
