.PHONY: build run test clean

build:
	@echo "building the application..."
	@go build -o ./bin/sso ./cmd/sso/main.go

run: build
	@echo "running the application..."
	@./bin/sso --config=./config/local.yaml

test:
	@echo "running tests..."
	@go test ./...

clean:
	@echo "cleaning up..."
	@rm -rf ./bin
