.PHONY: up build test help default

default: build

help:
	@echo 'Management commands for hellosrv:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make up              Start the service.'
	@echo '    make test-build      Build the project for testing.'
	@echo '    make test            Run tests and lint on a compiled project.'
	@echo '    make test-shell      Open shell with test environment.'
	@echo '    make test            Run go vet on a compiled project.'
	@echo

build:
	docker-compose build

up: 
	docker-compose up

# export PATH="$PATH:$(go env GOPATH)/bin"
proto-go:
	protoc -I . \
    --go_out ./ --go_opt paths=source_relative \
    --go-grpc_out ./ --go-grpc_opt paths=source_relative \
    api/**/**/*.proto