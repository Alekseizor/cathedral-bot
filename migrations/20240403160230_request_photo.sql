-- +goose Up
-- +goose StatementBegin
CREATE TABLE request_photo
(
    id            serial PRIMARY KEY,
    year          INT,
    study_program VARCHAR(100),
    event         VARCHAR(100),
    is_event_new  boolean      DEFAULT FALSE,
    description   VARCHAR(255),
    count_people  INT,
    marked_person INT,
    marked_people TEXT[],
    teachers      TEXT[],
    pointer       INT DEFAULT 0,
    attachment    VARCHAR(255),
    attachments   TEXT[],
    user_id       INT NOT NULL,
    status        INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS request_photo;
-- +goose StatementEnd