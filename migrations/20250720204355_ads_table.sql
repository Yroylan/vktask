-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ads (
 	id SERIAL PRIMARY KEY,
 	title TEXT NOT NULL,
 	description TEXT NOT NULL,
 	image_url TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 	price DECIMAL(10, 2) NOT NULL,
 	user_id INTEGER NOT NULL REFERENCES vkusers(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ads;
-- +goose StatementEnd
