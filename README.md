## GraphQL API Starter

### Setup

```bash
# install required tools.
$ bash ./scripts/tools.sh

# generate code.
$ task generate
```


### Generate JWT Secret

```bash
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