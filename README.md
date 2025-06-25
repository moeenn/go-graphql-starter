## GraphQL API Starter

A starter-kit go Golang projects using GraphQL APIs and PostgreSQL as database. Following tools and packages are used in this starter-kit:

- [gqlgen](https://gqlgen.com): GraphQL server implementation (utilizing code-generation)
- [sqlc](https://sqlc.dev): Library for accessing database (utilizing code-generation)
- [taskfile](https://taskfile.dev): CLI tool for running project scripts 


### Setup

```bash
# 01 - install required tools.
$ bash ./scripts/tools.sh

# 02 - create .env file.
$ cp .env.example .env

# 03 - generate code.
$ task generate

# 04 - generate the JWT secret token.
$ task gen:token
```

**Note**: Save the generated token in `.env` file against the key `JWT_SECRET`.


### Running the project

```bash
# run in development mode.
$ task build:run

# generate production build.
$ task build
```