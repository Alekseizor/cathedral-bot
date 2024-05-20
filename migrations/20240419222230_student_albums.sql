-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_albums
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    description   VARCHAR(10000),
    url           VARCHAR(255),
    vk_id         VARCHAR(255)
);

INSERT INTO student_albums (year, study_program, event, description, url, vk_id)
VALUES (2024, 'Бакалавриат', 'Защита диплома', 'Защита диплома бакалавриата',
        'https://vk.com/album-211704031_283523239', 283523239),
       (2024, 'Бакалавриат', 'Тазы', 'Мероприятие Тазы бакалавриата', 'https://vk.com/album-211704031_283523140',
        283523140),
       (2024, 'Бакалавриат', 'Учёба', 'Учёба бакалавриата', 'https://vk.com/album-211704031_283523140', 283523140),
       (2024, 'Магистратура', 'Защита диплома', 'Защита диплома магистратуры',
        'https://vk.com/album-211704031_283523141', 283523141),
       (2024, 'Магистратура', 'Выпускной', 'Выпускной магистратуры', 'https://vk.com/album-211704031_283523142',
        283523142),
       (2024, 'Специалитет', 'Тазы', 'Мероприятие Тазы специалитета', 'https://vk.com/album-211704031_283523142',
        283523142),
       (2024, 'Специалитет', 'Защита диплома', 'Защита диплома специалитета',
        'https://vk.com/album-211704031_283523142', 283523142),
       (2024, 'Аспирантура', 'Отдых', 'Отдых аспирантуры', 'https://vk.com/album-211704031_283523148', 283523148),
       (2024, 'Аспирантура', 'Защита диссертации', 'Защита диссертации аспирантуры',
        'https://vk.com/album-211704031_283523148', 283523148),

       (2023, 'Бакалавриат', 'Защита диплома', 'Защита диплома бакалавриата',
        'https://vk.com/album-211704031_283523239', 283523239),
       (2023, 'Бакалавриат', 'Тазы', 'Мероприятие Тазы бакалавриата', 'https://vk.com/album-211704031_283523140',
        283523140),
       (2023, 'Бакалавриат', 'Учёба', 'Учёба бакалавриата', 'https://vk.com/album-211704031_283523140', 283523140),
       (2023, 'Магистратура', 'Защита диплома', 'Защита диплома магистратуры',
        'https://vk.com/album-211704031_283523141', 283523141),
       (2023, 'Магистратура', 'Выпускной', 'Выпускной магистратуры', 'https://vk.com/album-211704031_283523142',
        283523142),
       (2023, 'Магистратура', 'Учёба', 'Учёба магистратуры', 'https://vk.com/album-211704031_283523142', 283523142),
       (2023, 'Магистратура', 'Отдых', 'Отдых магистратуры', 'https://vk.com/album-211704031_283523142', 283523142),
       (2023, 'Аспирантура', 'Защита диссертации', 'Защита диссертации аспирантуры',
        'https://vk.com/album-211704031_283523148', 283523148),

       (2022, 'Бакалавриат', 'Защита диплома', 'Защита диплома бакалавриата',
        'https://vk.com/album-211704031_283523239', 283523239),
       (2022, 'Бакалавриат', 'Тазы', 'Мероприятие Тазы бакалавриата', 'https://vk.com/album-211704031_283523140',
        283523140),
       (2022, 'Бакалавриат', 'Учёба', 'Учёба бакалавриата', 'https://vk.com/album-211704031_283523140', 283523140),
       (2022, 'Бакалавриат', 'Отдых', 'Отдых бакалавриата', 'https://vk.com/album-211704031_283523140', 283523140),
       (2022, 'Магистратура', 'Выпускной', 'Выпускной магистратуры', 'https://vk.com/album-211704031_283523142',
        283523142),
       (2022, 'Специалитет', 'Тазы', 'Мероприятие Тазы специалитета', 'https://vk.com/album-211704031_283523142',
        283523142),
       (2022, 'Специалитет', 'Защита диплома', 'Защита диплома специалитета',
        'https://vk.com/album-211704031_283523142', 283523142),
       (2022, 'Специалитет', 'Выпускной', 'Выпускной специалитета', 'https://vk.com/album-211704031_283523142',
        283523142),

       (2021, 'Бакалавриат', 'Защита диплома', 'Защита диплома бакалавриата',
        'https://vk.com/album-211704031_283523239', 283523239),
       (2021, 'Магистратура', 'Выпускной', 'Выпускной магистратуры', 'https://vk.com/album-211704031_283523142',
        283523142),

       (2020, 'Специалитет', 'Защита диплома', 'Защита диплома специалитета',
        'https://vk.com/album-211704031_283523142', 283523142),
       (2020, 'Аспирантура', 'Защита диссертации', 'Защита диссертации аспирантуры',
        'https://vk.com/album-211704031_283523148', 283523148),

       (2019, 'Бакалавриат', 'Учёба', 'Учёба бакалавриата', 'https://vk.com/album-211704031_283523140', 283523140);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_albums;
-- +goose StatementEnd