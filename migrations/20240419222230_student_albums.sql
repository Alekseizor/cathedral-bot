-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_albums
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    description         VARCHAR(10000),
    url           VARCHAR(255)
);

INSERT INTO student_albums (year, study_program, event,description, url)
VALUES (2024, 'Бакалавриат', 'Диплом','Получение диплома', 'https://vk.com/album-211704031_283523239'),
       (2024, 'Бакалавриат', 'Тазы', 'Веселье на тазах','https://vk.com/album-211704031_283523140'),
       (2024, 'Магистратура', 'Диплом','Мы магистры!', 'https://vk.com/album-211704031_283523141'),
       (2024, 'Магистратура', 'Выпускной', 'Выпускной, который невозможно забыть','https://vk.com/album-211704031_283523142'),

       (2023, 'Аспирантура', 'Выпускной','Выпускной кандидатов наук', 'https://vk.com/album-211704032_283523143'),
       (2023, 'Специалитет', 'Учёба','Век живи -век учись', 'https://vk.com/album-211704032_283523144'),

       (2022, 'Магистратура', 'Выпускной','Выпускной выезд', 'https://vk.com/album-211704033_283523145');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_albums;
-- +goose StatementEnd