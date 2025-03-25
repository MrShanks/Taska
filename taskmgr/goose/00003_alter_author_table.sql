-- +goose Up
-- +goose StatementBegin

ALTER TABLE author ALTER COLUMN password SET NOT NULL;
ALTER TABLE author ADD COLUMN token TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE author DROP COLUMN token;
ALTER TABLE author ALTER COLUMN password DROP NOT NULL;

-- +goose StatementEnd
