-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
    id         UUID PRIMARY KEY,
    title      VARCHAR(255),
    author     VARCHAR(100),
    year       INT,
    category   VARCHAR(50),
    hashtags   TEXT[],
    url        VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
