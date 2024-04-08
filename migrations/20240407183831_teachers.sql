-- +goose Up
-- +goose StatementBegin
CREATE TABLE teachers
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) unique
);
INSERT INTO teachers(name)
VALUES ('Филлипович Анна Юрьевна'),
       ('Канев Антон Иstateгоревич'),
       ('Гапанюк Юрий Евгеньевич'),
       ('Терехов Валерий Игоревич'),
       ('Черненький Михаил Валерьевич'),
       ('Лосева Светлана Сергеевна'),
       ('Спиридонов Сергей Борисович');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teachers;
-- +goose StatementEnd