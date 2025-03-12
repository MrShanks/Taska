-- +goose Up
-- +goose StatementBegin
create table tasks (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    title text UNIQUE,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tasks;
-- +goose StatementEnd
