-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS rates (
    timestamp TIMESTAMP NOT NULL,
    market TEXT NOT NULL,
    ask NUMERIC NOT NULL,
    bid NUMERIC NOT NULL 
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rates;
-- +goose StatementEnd
