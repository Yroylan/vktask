-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vkusers (
    id SERIAL PRIMARY KEY,
	login TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vkusers;
-- +goose StatementEnd
