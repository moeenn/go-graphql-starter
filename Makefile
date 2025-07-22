include .env

ENV_FILE = .env
ENTRYPOINT = ./cmd/api/main.go


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test ./...


gen_graphql:
	gqlgen generate


.PHONY: gen_token
gen_token:
	go-token


migration_new:
	goose -s create ${NAME} sql


.PHONY: db_migrate
db_migrate:
	goose up


.PHONY: db_rollback
db_rollback:
	goose down-to 0


.PHONY: run
run:
	godotenv -f ${ENV_FILE} go run ${ENTRYPOINT}


build: gen_graphql lint
	go build -o ./bin/api ${ENTRYPOINT}
