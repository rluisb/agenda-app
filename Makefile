build:
	@go build -o /build/main src/main.go

run: build
	@./build/main

test:
	@go test -v ./...

deploy-k8s:
	@kubectl apply -f k8s/mongodb-svc.yml
	@kubectl apply -f k8s/mongodb-statefulset.yml
	@kubectl apply -f k8s/Deployment.yml
	@kubectl apply -f k8s/Service.yml