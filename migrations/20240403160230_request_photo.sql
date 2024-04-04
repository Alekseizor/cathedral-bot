-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_photo
(
    id              serial PRIMARY KEY,
    year            INT,
    study_program   VARCHAR(100),
    event           VARCHAR(100),
    is_event_new    boolean DEFAULT FALSE,
    description     VARCHAR(255),
    marked_people   VARCHAR(255),
    attachment      VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_photo;
-- +goose StatementEnd