# Kubernetes IoT addon
IoT addon for Kubernetes clusters.

## Status
[![Build Status](https://travis-ci.org/fest-research/iot-addon.svg?branch=master)](https://travis-ci.org/fest-research/iot-addon)

## Development
Clone repository into `$GOPATH/src/github.com/fest-research/`:
```
mkdir -p $GOPATH/src/github.com/fest-research/
cd $GOPATH/src/github.com/fest-research/
git clone git@github.com:fest-research/iot-addon.git
```

To format source files use `govendor`, it will skip dependencies:

```
govendor fmt +local
```

## Usage
Use following commands to start all modules:

```
go run cmd/apiserver/apiserver.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
go run cmd/controller/controller.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
```

## Building Docker images
To build docker images use following command:
```
make build
```

To deploy it to Docker Hub use following commands:
```
docker login
make deploy
```
