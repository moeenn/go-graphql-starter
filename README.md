## GraphQL API Starter

A starter-kit go Golang projects using GraphQL APIs and PostgreSQL as database. Following tools and packages are used in this starter-kit:

- [gqlgen](https://gqlgen.com): GraphQL server implementation (utilizing code-generation)
- [sqlx](https://jmoiron.github.io/sqlx): Library for accessing database


### Setup

```bash
# 01 - install required tools.
$ bash ./scripts/tools.sh

# 02 - create .env file.
$ cp .env.example .env

# 03 - generate code.
$ make gen_graphql

# 04 - generate the JWT secret token.
$ make gen_token
```

**Note**: Save the generated token in `.env` file against the key `JWT_SECRET`.


### Running the project

```bash
# run in development mode.
$ make run

# generate production build.
$ make build
```
