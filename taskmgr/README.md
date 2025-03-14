# Postgres Database
Run the Database

using postgres 16.8-1.pgdg24.04+1


-- export databases variables
export POSTGRES_PWD="secure_password"


-- run the next commands to create a database and the user to manage it

sudo -i -u postgres psql
create database taskadb;
create user taskauser with superuser encrypted password '$POSTGRES_PWD';
grant all on database taskadb to taskauser;