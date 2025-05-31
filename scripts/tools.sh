#! /bin/bash

# install external tools.
go install -v github.com/go-task/task/v3/cmd/task@latest;
go install -v github.com/nametake/golangci-lint-langserver@latest;
go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
go install -v github.com/joho/godotenv/cmd/godotenv@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# install dependencies.
go mod tidy
