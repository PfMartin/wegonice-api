# wegonice-api

API for vegan recipes

## Initial Database Setup

1. Create a `.env` file similar to the `.example.env`
1. Go to `./database` and run `docker compose up -d` to start the wegonice-db
1. Create a new user, which can be used to connect to the database

```zsh
make db-create-user`
```

1. Now you can connect to the authentication database with

```zsh
make db-connect
```

## Create swagger documentation

- Install swag with `make get-swag`
- Generate swagger documentation with `make docs`
- Documenation is accessible on the route `{basePath}/docs/index.html`
