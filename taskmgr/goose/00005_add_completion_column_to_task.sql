-- +goose Up
-- +goose StatementBegin

ALTER TABLE task ADD COLUMN complete boolean NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE task DROP COLUMN complete;

-- +goose StatementEnd
