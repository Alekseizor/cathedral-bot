-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_albums
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    url           VARCHAR(255)
);

INSERT INTO student_albums (year, study_program, event, url)
VALUES (2024, 'Бакалавриат', 'Диплом', 'https://vk.com/album-211704031_283523139'),
       (2024, 'Бакалавриат', 'Тазы', 'https://vk.com/album-211704031_283523140'),


       (2023, 'Аспирантура', 'Выпускной', 'https://vk.com/album-211704032_283523143'),
       (2023, 'Специалитет', 'Учёба', 'https://vk.com/album-211704032_283523144'),
       (2024, 'Магистратура', 'Диплом', 'https://vk.com/album-211704031_283523141'),
       (2024, 'Магистратура', 'Выпускной', 'https://vk.com/album-211704031_283523142'),
       (2022, 'Магистратура', 'Выпускной', 'https://vk.com/album-211704033_283523145');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_albums;
-- +goose StatementEnd