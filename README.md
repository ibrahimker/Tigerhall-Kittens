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

To run integration test, we need to start all dependencies needed. We provide all dependencies via [Docker Compose](https://docs.docker.com/compose/).
Make sure to install [Docker Compose](https://docs.docker.com/compose/install/) before running integration test.

Also, we need to build the docker image for tigerhall-kittens first.

```sh
$ make compile-server
$ make docker-build-server
```

After that, run all images needed using `docker-compose`.

```sh
$ docker-compose up
```

Now that the server is run, you can hit it from postman / any http client