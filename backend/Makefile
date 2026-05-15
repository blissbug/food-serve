include .env
export

DB_URL = mysql://$(DBUser):$(DBPass)@tcp(localhost:3356)/$(DBName)?multiStatements=true

run:
	@go run cmd/api/main.go

migration:
	migrate create -ext sql -dir ./internal/db/migrations -seq create_$(name)_table

migrate-up:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" down

migrate-reset:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" force 0

db-reset:
	docker compose down -v
	docker compose up -d database redis
	sleep 5
	make migrate-up