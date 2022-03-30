# Tigerhall Kittens

## Description

Monolithic repository services for Tigerhall Kittens services.

## Test

### Unit Test

```sh
$ make tidy
$ make cover
```

### API Test

To run API test, we need to start all dependencies needed. We provide all dependencies via [Docker Compose](https://docs.docker.com/compose/).
Make sure to install [Docker Compose](https://docs.docker.com/compose/install/) before running integration test.

Also, we need to build the docker image for tigerhall-kittens first.

```sh
$ make rebuild-server
```

After that, run all images needed using `docker-compose`.

```sh
$ docker-compose up
```

Now that the server is run, we need to migrate our local database first before running the program. Read [Database Migration](doc/DATABASE_MIGRATION.md) for detail information

```sh
$ make migrate-schema url="postgres://postgresuser:postgrespassword@localhost:5432/tigerhall"
$ make migrate url="postgres://postgresuser:postgrespassword@localhost:5432/tigerhall" module=sighting
```

After you run migration, now the API is ready to be hit