# IoT API Server
A replacement to the k8s API server with minimal coverage of the kubelet's basic dependencies.

## Run
Go to `cmd/apiserver` and run:
```
go run apiserver.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
```

# IoT Controller
Kubernetes controller for IoT orchestration.

## Setup
Make sure, that you have valid `$GOPATH` set.

Clone repository into `$GOPATH/src/github.com/fest-research/`:
```
mkdir -p $GOPATH/src/github.com/fest-research/
cd $GOPATH/src/github.com/fest-research/
git clone git@github.com:fest-research/iot-controller.git
```

Run application:
Go to `cmd/controller` and run:
```
go run controller.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
```

## Tools
To format source files use `govendor`, it will skip dependencies:

```
govendor fmt +local
```