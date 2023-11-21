build:
	@go build -o bin/api src/main.go

run: build
	@./bin/api

test:
	@go test -v ./...