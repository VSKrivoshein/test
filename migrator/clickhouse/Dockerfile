FROM migrate/migrate

COPY /migrations /migrations

# Переопределяем entrypoin базового образа, чтобы поднять миграции в одной команде CMD
ENTRYPOINT []

CMD migrate -database 'clickhouse://host.docker.internal:9000?x-multi-statement=true' -path ./migrations/ up