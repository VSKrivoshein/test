#Только для локальной разработки без запуска приложения в докере
include .env
export
export REDIS_HOST=localhost
export KAFKA_HOST=localhost
export DB_HOST=localhost
export REDIS_PORT=6379

install:
	go get ./...

start:
	docker-compose --profile migrator up --build

evans:
	evans internal/app/api/grpc_api/proto/user.proto -p 8080

down:
	docker-compose down

proto:
	protoc -I=internal/app/api/grpc_api/proto --go-grpc_out=internal/app/api/grpc_api/proto --go_out=internal/app/api/grpc_api/proto internal/app/api/grpc_api/proto/*.proto

build:
	go build ./cmd/apiserver/main.go && ./main

