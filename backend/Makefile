.PHONY: run build migrate-up migrate-down

run:
	go run cmd/server/main.go

build:
	go build -o CargoLedger cmd/server/main.go

migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" down