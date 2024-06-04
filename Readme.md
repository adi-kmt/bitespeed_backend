## Bitespeed Backend Task

Implemented using go-fiber, postgres and sqlc.

First run

```shell
docker compose up
```

Make sure that [golang-migrate](https://github.com/golang-migrate/migrate) is installed.

Then run

```shell
make db-migrate
```
