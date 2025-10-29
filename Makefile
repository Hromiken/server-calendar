APP_NAME=server-calendar

.PHONY: build run test tidy

build:
	go build -o $(APP_NAME) ./cmd

run:
	go run ./cmd

test:
	go test ./internal/... -v

tidy:
	go mod tidy
	goimports -w .

