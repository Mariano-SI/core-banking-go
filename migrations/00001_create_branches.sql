-- +goose Up
CREATE TABLE branches (
    id         uuid        PRIMARY KEY,
    code       char(4)     NOT NULL UNIQUE,
    name       text        NOT NULL,
    active     boolean     NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down

DROP TABLE branches;
