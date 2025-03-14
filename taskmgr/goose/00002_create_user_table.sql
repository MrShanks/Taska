-- +goose Up
-- +goose StatementBegin
CREATE TABLE author (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    firstname TEXT,
    lastname TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT
);

TRUNCATE TABLE task;

ALTER TABLE task ADD COLUMN author_id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE task ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES author(id) ON DELETE CASCADE;

INSERT INTO author (firstname, lastname, email, password) VALUES('marco', 'rossi', 'marco@rossi.com', 'password')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE task DROP CONSTRAINT fk_author;
ALTER TABLE task DROP COLUMN author_id;
DROP TABLE author;
-- +goose StatementEnd
