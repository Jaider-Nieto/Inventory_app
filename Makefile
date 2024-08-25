build:
	@go build -o bin/ecommerce-go cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecommerce-go