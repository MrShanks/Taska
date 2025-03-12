-- create database
create database taskadb;

-- create user with permission to create database
create user taskauser with superuser encrypted password '$POSTGRES_PWD';

-- grant all permissions
grant all on database taskadb to taskauser;