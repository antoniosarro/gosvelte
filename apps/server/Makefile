include .env

create_migration:
	migrate create -ext=sql -dir=internal/migrations -seq init

migrate_up:
	migrate -path=internal/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

swagger:
	swag init -g cmd/main.go

.PHONY: create_migration migrate_up migrate_down swagger