-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_document_publication
(
    id            serial PRIMARY KEY,
    pointer       INT DEFAULT 0,
    user_id       INT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_document_publication;
-- +goose StatementEnd