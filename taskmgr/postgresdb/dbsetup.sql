-- create database
create database taskadb;

-- create user with permission to create database
create user taskauser with superuser encrypted password 'secure_password';

-- grant all permissions
grant all on database taskadb to taskauser;

-- connect to taskadb
\connect taskadb

-- create table tasks
create table tasks (id uuid DEFAULT gen_random_uuid() PRIMARY KEY,title text UNIQUE,description text);