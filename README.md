# :robot: Sensor Mockery server
A server that makes use of the sensor mockery lib to send mocked data over http

## How can I use this?
Currently you have two alternatives:

### Makefile
#### Prerequisites
You need to have the following tools installed:
- :mouse2: go
- :whale: docker
- direnv
- make
Also you would need a Unix (or Unix-like) OS. The following environment variables need to be set - GOBIN, GOPATH.

#### Use cases
The following make commands are supported:
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
```
