-- +goose Up
-- +goose StatementBegin
CREATE TABLE view_request_photo
(
    id            serial PRIMARY KEY,
    pointer       INT DEFAULT 0,
    user_id       INT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS view_request_photo;
-- +goose StatementEnd