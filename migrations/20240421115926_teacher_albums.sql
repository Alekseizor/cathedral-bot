-- +goose Up
-- +goose StatementBegin
CREATE TABLE teacher_albums
(
    id      serial PRIMARY KEY,
    name VARCHAR(100),
    url     VARCHAR(255)
);

INSERT INTO teacher_albums (name, url)
VALUES ('Абулкасимов Манас Мукитович', 'https://vk.com/album-211704031_283523239'),
       ('Аксёнова Мария Владимировна', 'https://vk.com/album-211704031_283523240'),
       ('Аладин Дмитрий Владимирович', 'https://vk.com/album-211704031_283523241'),
       ('Али Диб Ахмада', 'https://vk.com/album-211704031_283523242'),
       ('Антонов Артем Ильич', 'https://vk.com/album-211704031_283523243'),
       ('Афанасьев Арсений Геннадьевич', 'https://vk.com/album-211704031_283523244'),
       ('Афанасьев Геннадий Иванович', 'https://vk.com/album-211704031_283523245'),
       ('Балашов Антон Михайлович', 'https://vk.com/album-211704031_283523246'),
       ('Балдин Александр Викторович', 'https://vk.com/album-211704031_283523247'),
       ('Белодедов Михаил Владимирович', 'https://vk.com/album-211704031_283523248'),
       ('Большаков Сергей Алексеевич', 'https://vk.com/album-211704031_283523249'),
       ('Булатова Ирина Георгиевна', 'https://vk.com/album-211704031_283523250'),
       ('Варламов Олег Олегович', 'https://vk.com/album-211704031_283523251'),
       ('Виноградова Мария Валерьевна', 'https://vk.com/album-211704031_283523252'),
       ('Волков Артём Сергеевич', 'https://vk.com/album-211704031_283523253'),
       ('Галкин Валерий Александрович', 'https://vk.com/album-211704031_283523254'),
       ('Гапанюк Юрий Евгеньевич', 'https://vk.com/album-211704031_283523255'),
       ('Горячкин Борис Сергеевич', 'https://vk.com/album-211704031_283523256'),
       ('Григорьев Юрий Александрович', 'https://vk.com/album-211704031_283523257'),
       ('Дятленко Елена Александровна', 'https://vk.com/album-211704031_283523258'),
       ('Ишков Денис Олегович', 'https://vk.com/album-211704031_283523259'),
       ('Калистратов Алексей Павлович', 'https://vk.com/album-211704031_283523260'),
       ('Канев Антон Игоревич', 'https://vk.com/album-211704031_283523261'),
       ('Карабулатова Ирина Советовна', 'https://vk.com/album-211704031_283523262'),
       ('Кесель Сергей Александрович', 'https://vk.com/album-211704031_283523263'),
       ('Ковалева Наталья Александровна', 'https://vk.com/album-211704031_283523264'),
       ('Кротов Юрий Николаевич', 'https://vk.com/album-211704031_283523265'),
       ('Крутов Тимофей Юрьевич', 'https://vk.com/album-211704031_283523266'),
       ('Лабунец Леонид Витальевич', 'https://vk.com/album-211704031_283523267'),
       ('Лосева Светлана Сергеевна', 'https://vk.com/album-211704031_283523268'),
       ('Максаков Алексей Александрович', 'https://vk.com/album-211704031_283523269'),
       ('Масленников Константин Юрьевич', 'https://vk.com/album-211704031_283523270'),
       ('Машкин Константин Вилиорович', 'https://vk.com/album-211704031_283523271'),
       ('Михеев Вячеслав Алексеевич', 'https://vk.com/album-211704031_283523272'),
       ('Мышенков Константин Сергеевич', 'https://vk.com/album-211704031_283523273'),
       ('Нардид Анатолий Николаевич', 'https://vk.com/album-211704031_283523274'),
       ('Нестеров Юрий Григорьевич', 'https://vk.com/album-211704031_283523275'),
       ('Плужникова Ольга Юрьевна', 'https://vk.com/album-211704031_283523276'),
       ('Постников Виталий Михайлович', 'https://vk.com/album-211704031_283523277'),
       ('Правдина Анна Дмитриевна', 'https://vk.com/album-211704031_283523278'),
       ('Попов Илья Андреевич', 'https://vk.com/album-211704031_283523279'),
       ('Самохвалов Алексей Эдуардович', 'https://vk.com/album-211704031_283523280'),
       ('Селивёрстова Анастасия Валерьевна', 'https://vk.com/album-211704031_283523281'),
       ('Семёнов Дмитрий Валериевич', 'https://vk.com/album-211704031_283523282'),
       ('Семкин Петр Степанович', 'https://vk.com/album-211704031_283523283'),
       ('Силантьева Елена Юрьевна', 'https://vk.com/album-211704031_283523284'),
       ('Симонов Михаил Фёдорович', 'https://vk.com/album-211704031_283523285'),
       ('Спиридонов Сергей Борисович', 'https://vk.com/album-211704031_283523286'),
       ('Строганов Виктор Юрьевич', 'https://vk.com/album-211704031_283523287'),
       ('Строганов Дмитрий Викторович', 'https://vk.com/album-211704031_283523288'),
       ('Сухобоков Андрей Валентинович', 'https://vk.com/album-211704031_283523289'),
       ('Терехов Валерий Игоревич', 'https://vk.com/album-211704031_283523290'),
       ('Тимофеев Виктор Борисович', 'https://vk.com/album-211704031_283523291'),
       ('Филиппович Анна Юрьевна', 'https://vk.com/album-211704031_283523292'),
       ('Хайруллин Рустам Зиннаттулович', 'https://vk.com/album-211704031_283523293'),
       ('Черненький Михаил Валерьевич', 'https://vk.com/album-211704031_283523294'),
       ('Черненький Станислав Валерьевич', 'https://vk.com/album-211704031_283523295'),
       ('Шкуратова Людмила Петровна', 'https://vk.com/album-211704031_283523296'),
       ('Шук Владимир Павлович', 'https://vk.com/album-211704031_283523297'),
       ('Якубов Алексей Ренатович', 'https://vk.com/album-211704031_283523298');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teacher_albums;
-- +goose StatementEnd