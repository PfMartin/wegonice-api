# wegonice-api

API for vegan recipes

## Initial Database Setup

- Create a `.env` file similar to the `.example.env`
- Go to `./database` and run `docker compose up -d` to start the wegonice-db
- Create a new user, which can be used to connect to the database

```zsh
make db-create-user`
```

- Now you can connect to the authentication database with

```zsh
make db-connect
```

## Server

- Start the server with

```zsh
go run .
```

- Test the server with

```zsh
go test ./...
go test -race ./... -run 'Unit'
```

## Create swagger documentation

- Install swag with `make get-swag`
- Generate swagger documentation with `make docs`
- Documenation is accessible on the route `{basePath}/docs/index.html`
