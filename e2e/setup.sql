-- Insert authors
INSERT INTO author (id, firstname, lastname, email, password) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Marco', 'Rossi', 'marco@rossi.com', 'password'),
    ('22222222-2222-2222-2222-222222222222', 'Luca', 'Bianchi', 'luca@bianchi.com', 'password');

-- Insert tasks
INSERT INTO task (id, title, description, author_id) VALUES
    ('33333333-3333-3333-3333-333333333333', 'Task 1', 'Description for Task 1', '11111111-1111-1111-1111-111111111111'),
    ('44444444-4444-4444-4444-444444444444', 'Task 2', 'Description for Task 2', '11111111-1111-1111-1111-111111111111'),
    ('55555555-5555-5555-5555-555555555555', 'Task 3', 'Description for Task 3', '22222222-2222-2222-2222-222222222222');
