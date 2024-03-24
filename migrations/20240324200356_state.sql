-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS state (
    "vk_id" integer not null primary key,
    "title" varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS state;
-- +goose StatementEnd
