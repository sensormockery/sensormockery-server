# Development guide

## Prerequisites

You need to have the following tools installed:
- :mouse2: go
- :whale: docker
- direnv
- make
- [golint](https://github.com/golang/lint)

Also it would be easier to be on a Unix (or Unix-like) OS. The following environment variables need to be set - GOBIN, GOPATH.

For the sake of simplicity we are using only `postgre` as our db.

## :wrench: Workflow

We are following a standard workflow:
- each new feature is developed in a `feature branch`
- there is no `dev branch` and everything is pushed merged directly into `master`
- run all the tests with at least `make run-system-tests` before making a PR
- good to request reviews
- hacky and untested scripts should go to the `hack` folder in the root of the project

Before committing:
- run `make clean-code`
- run `shellcheck` on scripts

## Technologies

The server is written in `go` and tests are respectively in `ginkgo`. For db migrations we are using [go-migrate](https://github.com/golang-migrate/migrate). Helper scripts are written in `GNU Makefile`.

## Makefile

A lot of the work is done through `Makefile`. The following make commands are supported:

```bash
make # installs a local binary to your GOBIN and builds a docker image
make build # the same as make
make clean # remove local binary from GOBIN and delete local docker image

make check-compliance # checks if you have the prerequisites set

make build-binary # installs a local binary to your GOBIN
make clean-binary # removes your local binary from GOBIN

make build-docker-image # builds a docker image
make clean-docker-image # delete local docker image
make run-docker-image # run docker image
make push-docker-image # push docker image to a docker registry

# testing
make run-system-tests
make run-system-tests skip_update=true # to skip rebuilding when running tests
```

## Environment

Most of the variable you need would be loaded from `.envrc`. You would additionally need:
- `DOCKERHUB_USERNAME` - Your dockerhub username
- `DOCKERHUB_PASSWORD` - Your dockerhub password
- `DB_USER` - Your postgre db user
- `DB_PASS` - Your postgre db password
