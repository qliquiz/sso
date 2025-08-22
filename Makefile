.PHONY: build run test clean migrate-up migrate-down migrate-tests-up migrate-tests-down

build:
	@echo "building the application..."
	@go build -o ./bin/sso ./cmd/sso/main.go

run: build
	@echo "running the application..."
	@./bin/sso --config=./config/local.yaml

test:
	@echo "running tests..."
	@go test ./tests

clean:
	@echo "cleaning up..."
	@rm -rf ./bin

# apply all migrations
migrate-up:
	@go run ./cmd/migrator \
		--config-path=./config/local.yaml \
		--migrations-path=./migrations \
		--command=up

# rollback the last migration
migrate-down:
	@go run ./cmd/migrator \
		--config-path=./config/local.yaml \
		--migrations-path=./migrations \
		--command=down

migrate-tests-up:
	@go run ./cmd/migrator \
		--config-path=./config/test.yaml \
		--migrations-path=./tests/migrations \
		--command=up

migrate-tests-down:
	@go run ./cmd/migrator \
		--config-path=./config/test.yaml \
		--migrations-path=./tests/migrations \
		--command=down
