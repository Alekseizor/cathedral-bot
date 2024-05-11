-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) unique
);

INSERT INTO events(name)
VALUES ('Учёба'),
       ('Защита диплома'),
       ('Защита диссертации'),
       ('Выпускной'),
       ('Тазы'),
       ('Отдых');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd