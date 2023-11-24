build:
	@go build -o /build/main src/main.go

run: build
	@./build/main

test:
	@go test -v ./...

deploy-k8s:
	@kubectl apply -f k8s

deploy-k8s-destroy:
	@kubectl delete -f k8s