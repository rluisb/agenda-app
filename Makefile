build:
	@go build -o bin/api src/main.go

run: build
	@./bin/api

test:
	@go test -v ./...

deploy-k8s:
	@kubectl apply -f k8s/namespace.yml
	@kubectl apply -f k8s/deployment.yml
	@kubectl apply -f k8s/service.yml