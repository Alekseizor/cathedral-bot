-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
                            id   SERIAL PRIMARY KEY,
                            name VARCHAR(255) unique
);
INSERT INTO events(name) VALUES ('Учёба');
INSERT INTO events(name) VALUES ('Диплом');
INSERT INTO events(name) VALUES ('Выпускной');
INSERT INTO events(name) VALUES ('Тазы');
INSERT INTO events(name) VALUES ('Отдых');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd