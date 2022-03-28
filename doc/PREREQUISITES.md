## Prerequisites

- Adjust Git Conifg

    By default, go will use `https` prefix to download dependencies. In case you found difficulties when downloading from private repo, adjust `~/.gitconfig`
    ```
    [url "ssh://git@github.com/"]
      insteadOf = https://github.com/
    ```

- Install Go

    We use version 1.17. Follow [Golang installation guideline](https://golang.org/doc/install).

- Install golangci-lint

    Follow [golangci-lint installation](https://golangci-lint.run/usage/install/).

- Install gomock

    Follow [gomock installation](https://github.com/golang/mock).

- Install golang-migrate

    Follow [golang-migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).