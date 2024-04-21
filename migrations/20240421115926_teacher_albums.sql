-- +goose Up
-- +goose StatementBegin
CREATE TABLE teacher_albums
(
    id               serial PRIMARY KEY,
    teacher          VARCHAR(100),
    url              VARCHAR(255)
);

INSERT INTO teacher_albums (teacher, url)
VALUES
    ('Абулкасимов Манас Мукитович', 'https://example.com/1'),
    ('Аксёнова Мария Владимировна', 'https://example.com/2'),
    ('Аладин Дмитрий Владимирович', 'https://example.com/3'),
    ('Али Диб Ахмада', 'https://example.com/4'),
    ('Антонов Артем Ильич', 'https://example.com/5'),
    ('Афанасьев Арсений Геннадьевич', 'https://example.com/6'),
    ('Афанасьев Геннадий Иванович', 'https://example.com/7'),
    ('Балашов Антон Михайлович', 'https://example.com/8'),
    ('Балдин Александр Викторович', 'https://example.com/9'),
    ('Белодедов Михаил Владимирович', 'https://example.com/10'),
    ('Большаков Сергей Алексеевич', 'https://example.com/11'),
    ('Булатова Ирина Георгиевна', 'https://example.com/12'),
    ('Варламов Олег Олегович', 'https://example.com/13'),
    ('Виноградова Мария Валерьевна', 'https://example.com/14'),
    ('Волков Артём Сергеевич', 'https://example.com/15'),
    ('Галкин Валерий Александрович', 'https://example.com/16'),
    ('Гапанюк Юрий Евгеньевич', 'https://example.com/17'),
    ('Горячкин Борис Сергеевич', 'https://example.com/18'),
    ('Григорьев Юрий Александрович', 'https://example.com/19'),
    ('Дятленко Елена Александровна', 'https://example.com/20'),
    ('Ишков Денис Олегович', 'https://example.com/21'),
    ('Калистратов Алексей Павлович', 'https://example.com/22'),
    ('Канев Антон Игоревич', 'https://example.com/23'),
    ('Карабулатова Ирина Советовна', 'https://example.com/24'),
    ('Кесель Сергей Александрович', 'https://example.com/25'),
    ('Ковалева Наталья Александровна', 'https://example.com/26'),
    ('Кротов Юрий Николаевич', 'https://example.com/27'),
    ('Крутов Тимофей Юрьевич', 'https://example.com/28'),
    ('Лабунец Леонид Витальевич', 'https://example.com/29'),
    ('Лосева Светлана Сергеевна', 'https://example.com/30'),
    ('Максаков Алексей Александрович', 'https://example.com/31'),
    ('Масленников Константин Юрьевич', 'https://example.com/32'),
    ('Машкин Константин Вилиорович', 'https://example.com/33'),
    ('Михеев Вячеслав Алексеевич', 'https://example.com/34'),
    ('Мышенков Константин Сергеевич', 'https://example.com/35'),
    ('Нардид Анатолий Николаевич', 'https://example.com/37'),
    ('Нестеров Юрий Григорьевич', 'https://example.com/38'),
    ('Плужникова Ольга Юрьевна', 'https://example.com/39'),
    ('Постников Виталий Михайлович', 'https://example.com/40'),
    ('Правдина Анна Дмитриевна', 'https://example.com/41'),
    ('Попов Илья Андреевич', 'https://example.com/42'),
    ('Самохвалов Алексей Эдуардович', 'https://example.com/43'),
    ('Селивёрстова Анастасия Валерьевна', 'https://example.com/44'),
    ('Семёнов Дмитрий Валериевич', 'https://example.com/45'),
    ('Семкин Петр Степанович', 'https://example.com/46'),
    ('Силантьева Елена Юрьевна', 'https://example.com/47'),
    ('Симонов Михаил Фёдорович', 'https://example.com/48'),
    ('Спиридонов Сергей Борисович', 'https://example.com/49'),
    ('Строганов Виктор Юрьевич', 'https://example.com/50'),
    ('Строганов Дмитрий Викторович', 'https://example.com/51'),
    ('Сухобоков Андрей Валентинович', 'https://example.com/52'),
    ('Терехов Валерий Игоревич', 'https://example.com/53'),
    ('Тимофеев Виктор Борисович', 'https://example.com/54'),
    ('Филиппович Анна Юрьевна', 'https://example.com/55'),
    ('Хайруллин Рустам Зиннаттулович', 'https://example.com/56'),
    ('Черненький Михаил Валерьевич', 'https://example.com/57'),
    ('Черненький Станислав Валерьевич', 'https://example.com/58'),
    ('Шкуратова Людмила Петровна', 'https://example.com/59'),
    ('Шук Владимир Павлович', 'https://example.com/60'),
    ('Якубов Алексей Ренатович', 'https://example.com/61');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teacher_albums;
-- +goose StatementEnd