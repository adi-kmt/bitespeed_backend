## Bitespeed Backend Task

Implemented using go-fiber, postgres and sqlc.

### Application Structure

- The first `db/` contains all the SQL queries with `/migration/` including the db migration
and `/queries/` having all the required SQL queries used in this application.

- `cmd/` contains the entry point of the application

- `pkg/` contains the application code which is further divided into `/controllers/`, `/services/` and 
`/repositories/` as per the MVC pattern.
    - Some other packages are used for utilities like finding the unique elements in a slice, implementing a simplistic DI mechanism, etc

First run

```shell
docker compose up
```

Make sure that [golang-migrate](https://github.com/golang-migrate/migrate) is installed.

Then run

```shell
make db-migrate
```

### Endpoints

```shell
/identify
```

POST Request with body 
```json
{
    "email": "email@sample.com",
    "phoneNumber": "8070695"
}
```

with response
```json
{
    "contact": {
        "primaryContactId": 1,
        "emails": [
            "sample@email.com",
            "email@sample.com "
        ],
        "phoneNumbers": [
            "8070695",
            "8070690"
        ],
        "secondaryContactNumbers": [
            2
        ]
    }
}
```
