-- +goose Up
CREATE SEQUENCE account_number_seq;

CREATE TABLE accounts (
    id          uuid        PRIMARY KEY,
    number      text        NOT NULL UNIQUE,
    branch_id   uuid        NOT NULL REFERENCES branches (id),
    customer_id uuid        NOT NULL UNIQUE REFERENCES customers (id),
    status      text        NOT NULL CHECK (status IN ('active', 'blocked', 'closed')),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE accounts;
DROP SEQUENCE account_number_seq;
