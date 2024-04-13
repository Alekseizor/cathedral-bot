-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_photo
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    is_event_new  boolean DEFAULT FALSE,
    description   VARCHAR(255),
    count_people  INT,
    marked_person INT,
    marked_people TEXT[],
    teachers      TEXT[],
    attachment    VARCHAR(255),
    user_id       INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_photo;
-- +goose StatementEnd