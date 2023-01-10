.PHONY: test
test:
	go test ./...
#
build: ## builds the project and binary
	go mod tidy
	go mod vendor
	go build -o benchmark ./tools/benchmark/main.go


.PHONY: test-all
test-all:
	docker-compose -f docker-compose.bdd.yaml build && docker-compose -f docker-compose.bdd.yaml run --rm test && docker-compose -f docker-compose.bdd.yaml down  --remove-orphans


.PHONY: lint
lint:
	golangci-lint run ./...