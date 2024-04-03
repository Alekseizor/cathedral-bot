-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admin (
    "vk_id" integer not null primary key,
    "documents" BOOLEAN,
    "albums" BOOLEAN
);
INSERT INTO admin VALUES (236322856,TRUE,true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS admin;
-- +goose StatementEnd
