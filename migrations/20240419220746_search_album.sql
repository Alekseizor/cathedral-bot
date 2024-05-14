-- +goose Up
-- +goose StatementBegin
CREATE TABLE search_album
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    surname       VARCHAR(100),
    pointer       INT DEFAULT 0,
    user_id       INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS search_album;
-- +goose StatementEnd