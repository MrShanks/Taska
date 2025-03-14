-- +goose Up
-- +goose StatementBegin
CREATE TABLE author (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    firstname TEXT,
    lastname TEXT,
    email TEXT UNIQUE NOT NULL,
    password TEXT
);

ALTER TABLE tasks ADD COLUMN author_id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE tasks ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES author(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP CONSTRAINT fk_author;
ALTER TABLE tasks DROP COLUMN author_id;
DROP TABLE author;
-- +goose StatementEnd
