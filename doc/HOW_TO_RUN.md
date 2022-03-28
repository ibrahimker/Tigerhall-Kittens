## How to Run

To run the application, we use docker images and docker compose.

### Docker

- Install [Docker Compose](https://docs.docker.com/compose/).

- Download the dependencies

    ```
    $ make tidy
    ```

- Compile the binary

    ```
    $ make compile-server
    ```

- Build image

    ```
    $ make docker-build-server
    ```

- Run docker compose

    ```
    $ docker-compose up
    ```