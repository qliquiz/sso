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

# apply all migrations
migrate-up:
	@go run ./cmd/migrator \
		--config-path=./config/storage.yaml \
		--migrations-path=./migrations \
		--command=up

# rollback the last migration
migrate-down:
	@go run ./cmd/migrator \
		--config-path=./config/storage.yaml \
		--migrations-path=./migrations \
		--command=down
