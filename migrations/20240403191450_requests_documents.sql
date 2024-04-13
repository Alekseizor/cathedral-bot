-- +goose Up
-- +goose StatementBegin
CREATE TABLE requests_documents
(
    id              serial PRIMARY KEY,
    title           VARCHAR(255),
    author          VARCHAR(100),
    year            INT,
    category        VARCHAR(50),
    is_category_new boolean DEFAULT FALSE,
    description     VARCHAR(255),
    hashtags        TEXT[],
    attachment      VARCHAR(255),
    user_id         INT NOT NULL,
    status          INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
