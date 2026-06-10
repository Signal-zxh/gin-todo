FROM mysql:8.0
COPY ./scripts/create-tables.sql /docker-entrypoint-initdb.d/