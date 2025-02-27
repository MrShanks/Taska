-- create database
create database taskadb;

-- create user with permission to create database
create user taskauser with superuser encrypted password 'secure_password';

-- grant all permissions
grant all on database taskadb to taskauser;

-- connect to taskadb
\connect taskadb

-- create table tasks
create table tasks (id uuid PRIMARY KEY,title text UNIQUE,description text);


insert into tasks (id, title , description) values ('4ed7d963-a74c-4231-92dd-f88e2e0b15f9', 'First task', 'first of all');
insert into tasks (id, title , description) values ('b5ff1aca-661d-47d5-8aae-4d0e99ee993f', 'Second task', 'and then');
insert into tasks (id, title , description) values ('17a9d5ef-c9bd-423d-b59f-3be1278a31c9', 'Third task', 'finally');