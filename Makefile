include .env
export

APP_NAME=be-technical-test
BUILD_DIR=./bin
MAIN_PATH=./cmd/main.go
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

tidy:
	go mod tidy

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME).exe $(MAIN_PATH)

run:
	@go run $(MAIN_PATH)

migrate-up:
	@migrate -path migration -database "$(DB_URL)" up
