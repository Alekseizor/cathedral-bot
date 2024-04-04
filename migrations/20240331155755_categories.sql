-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) unique
);
INSERT INTO Categories(Name) VALUES ('Статья');
INSERT INTO Categories(Name) VALUES ('Выпуск газеты');
INSERT INTO Categories(Name) VALUES ('Методическое указание');
INSERT INTO Categories(Name) VALUES ('Учебник');
INSERT INTO Categories(Name) VALUES ('Сборник учебно-методических работ и статей');
INSERT INTO Categories(Name) VALUES ('Указ');
INSERT INTO Categories(Name) VALUES ('Информация о преподавателях');
INSERT INTO Categories(Name) VALUES ('Отчет');
INSERT INTO Categories(Name) VALUES ('Курсовая работа');
INSERT INTO Categories(Name) VALUES ('Дипломная работа');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
