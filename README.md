# Sword Health API

Sword Health task management.

## Requirements

- [go](https://tip.golang.org/doc/go1.19)
- [mockgen](https://github.com/golang/mock)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [swaggo](https://github.com/swaggo/swag)

## Instalation

```bash
$ go get
```

## Configuration

### Creating .env file

For local environment just create a `.env` like this:

```txt
MIGRATE_URL=mysql://user:S3cR31@\(127.0.0.1:3306\)/swordhealth
MYSQL_PASSWORD=S3cR31
CRYPTO_HASH_KEY=FDJ1mnhuzjFjTdwhq7DtZG2Cq9kuuEZCG
CRYPTO_JWT_KEY=LmCD2mlHSQUiiu5AaYlXfrDOrnZaYzXCh
```

### Migrate

After running `docker-compose`, it's necessary to wait a few seconds to run the `migrate` command.

```bash
$ docker-compose up -d
$ make migrate
```

## Running

```bash
$ make run
```

See API local documentation at [swagger](http:localhost:8080/api/swagger/index.html)

---

## Local default manager user

To create other users, authenticate with manager user on `/api/auth/login` with BasicAuth:

```txt
Username = admin
Password = S3cR31 
```

Then, configure JWT authentication with access token response from login route, like this:

```txt
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjMyNDcxNjMsImlhdCI6MTY2MzI0NjI2Mywicm9sZSI6Im1hbmFnZXIiLCJzdWIiOjEsInVzZXJuYW1lIjoiYWRtaW4ifQ.vfZayJ5FQ15gyJKbQVSjgwNoeodV6ElICbJIfxz0buc
```

---

## Tests

```bash
# Unit tests
$ make test

# Tests coverage
$ make test/cov

# Tests benchmark service module
$ make test/bench
```

---

## Next steps

- Add update tasks flow
- Add delete tasks flow
- Add message broker to notify manager users
- Add deploy configuration files
