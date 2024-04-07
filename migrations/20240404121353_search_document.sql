-- +goose Up
-- +goose StatementBegin
CREATE TABLE search_document
(
    id          serial PRIMARY KEY,
    title       VARCHAR(255),
    author      VARCHAR(100),
    year        INT,
    start_year  INT,
    end_year    INT,
    categories  TEXT[],
    hashtags    TEXT[],
    documents   INT[],
    pointer_doc INT DEFAULT 0,
    user_id     INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
