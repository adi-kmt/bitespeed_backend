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

### Endpoints

```json
/identify

POST Request with body 
{
    "email": "email@sample.com ",
    "phoneNumber": "8070695"
}

with response

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
            1,
            2
        ]
    }
}
```
