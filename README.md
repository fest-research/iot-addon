# Kubernetes IoT addon
IoT addon for Kubernetes clusters.

## Setup
Make sure, that you have valid `$GOPATH` set.

Clone repository into `$GOPATH/src/github.com/fest-research/`:
```
mkdir -p $GOPATH/src/github.com/fest-research/
cd $GOPATH/src/github.com/fest-research/
git clone git@github.com:fest-research/iot-addon.git
```

## `apiserver`
IoT addon API server with minimal coverage of the `kubelet`'s basic dependencies.

Use following command to run:
```
go run cmd/apiserver/apiserver.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
```

## `controller`
IoT addon controller.

Use following command to run:
```
go run cmd/controller/controller.go --kubeconfig=<kubeconfig-path> --apiserver=<apiserver-adress>
```

## Tools
To format source files use `govendor`, it will skip dependencies:

```
govendor fmt +local
```