# Postgres Database
Run the Database

using postgres 16.8-1.pgdg24.04+1


-- create the db and user
sudo -i -u postgres psql -f taskmgr/postgresdb/dbsetup.sql


-- export databases variables
export POSTGRES_PWD="secure_password"