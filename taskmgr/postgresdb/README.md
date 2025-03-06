# Postgres Database
Run the Database

using postgres 16.8-1.pgdg24.04+1


-- start the db and create the tasks table
sudo -i -u postgres psql -f taskmgr/postgresdb/dbsetup.sql


-- export databases variables
export POSTGRES_USER="taskauser"
export POSTGRES_PWD="secure_password"
export POSTGRES_HOST="localhost"
export POSTGRES_PORT="5432"
export DB="taskadb"
