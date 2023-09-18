all: dependencies format imports test doc
dependencies:
	@echo "Syncing dependencies with go mod tidy"
	@go mod tidy
format:
	@echo "Formatting Go code recursively"
	@go fmt ./...
imports:
	@echo "Executing goimports recursively"
	@goimports -w $(find . -type f -name '*.go') ./
test:
	@echo "Running tests"
	@go test -v ./... -covermode=atomic -coverpkg=./... -count=1  -race -timeout=30m -shuffle=on
doc:
	@swag init --parseDependency --parseInternal -g "./cmd/api/main.go"
	@cp ./docs/swagger.yaml ./docs/specs