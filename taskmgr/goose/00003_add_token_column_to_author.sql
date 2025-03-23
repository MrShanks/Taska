
-- +goose Up
-- +goose StatementBegin

ALTER TABLE author ADD COLUMN token TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE author DROP COLUMN token;

-- +goose StatementEnd
