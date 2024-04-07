-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS object_admin (
    "admin_id" integer not null primary key,
    "object_id" integer
);
INSERT INTO object_admin VALUES (236322856);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS object_admin;
-- +goose StatementEnd
