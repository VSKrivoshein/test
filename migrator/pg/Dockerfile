FROM migrate/migrate

COPY /migrations /migrations

# Переопределяем entrypoin базового образа, чтобы поднять миграции в одной команде CMD
ENTRYPOINT []

CMD migrate -path=/migrations/ -database postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_TABLE}?sslmode=disable up