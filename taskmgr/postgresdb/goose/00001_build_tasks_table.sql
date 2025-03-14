-- +goose Up
-- +goose StatementBegin
create table task (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    title text UNIQUE NOT NULL,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table task;
-- +goose StatementEnd
