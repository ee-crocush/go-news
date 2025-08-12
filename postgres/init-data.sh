#!/bin/bash
set -e;


if [ -n "${$POSTGRES_USER:-}" ] && [ -n "${POSTGRES_PASSWORD:-}" ]; then
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
		CREATE DATABASE ${DB_NAME};
	EOSQL
else
	echo "SETUP INFO: No Environment variables given!"
fi