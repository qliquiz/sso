.PHONY: build run test clean

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
	@rm -rf ./bin ./api/gen

proto:
	@echo "generating proto files..."
	@mkdir api/gen
	@protoc \
		-I api/proto \
		api/proto/*.proto \
		--go_out=./api/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./api/gen \
		--go-grpc_opt=paths=source_relative
	@echo "proto files generated successfully."

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
