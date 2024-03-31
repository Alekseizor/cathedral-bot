-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255)
);
INSERT INTO Categories(Name) VALUES ('статья');
INSERT INTO Categories(Name) VALUES ('выпуск газеты');
INSERT INTO Categories(Name) VALUES ('методическое указание');
INSERT INTO Categories(Name) VALUES ('учебник');
INSERT INTO Categories(Name) VALUES ('сборник учебно-методических работ и статей');
INSERT INTO Categories(Name) VALUES ('указ');
INSERT INTO Categories(Name) VALUES ('информация о преподавателях');
INSERT INTO Categories(Name) VALUES ('отчет');
INSERT INTO Categories(Name) VALUES ('курсовая работа');
INSERT INTO Categories(Name) VALUES ('дипломная работа');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
