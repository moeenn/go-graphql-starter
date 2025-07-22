#! /bin/bash

# install external tools.
go install -v github.com/joho/godotenv/cmd/godotenv@latest
go install -v github.com/99designs/gqlgen@v0.17.73
go install -v github.com/pressly/goose/v3/cmd/goose@v3.24.3
go install -v github.com/moeenn/go-token@latest

# install dependencies.
go mod tidy
