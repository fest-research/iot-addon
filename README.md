# Kubernetes IoT addon
IoT addon for Kubernetes clusters.

## Status
[![Build Status](https://travis-ci.org/fest-research/iot-addon.svg?branch=master)](https://travis-ci.org/fest-research/iot-addon)

## Deploy to Kubernetes

```shell
$ kubectl create -f https://raw.githubusercontent.com/fest-research/iot-addon/master/assets/iot-addon.yaml
```
Flash RaspberryPi devices with this [software](https://github.com/fest-research/ubikube) to connect easily to the iot-server

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

## Configure kubectl

```shell
$ kubectl config set-cluster demo-cluster --server=http://104.155.11.172:8080/
$ kubectl config set-context demo --cluster=demo-cluster
$ kubectl config use-context demo
```

## Setup of Insecure Kubernetes
```
$ git clone https://github.com/kubernetes/kube-deploy
$ cd kube-deploy/docker-multinode
$ docker pull zreigz/hyperkube-amd64:v1.6.0-alpha.10
$ docker tag zreigz/hyperkube-amd64:v1.6.0-alpha.10 gcr.io/google_containers/hyperkube-amd64:v1.6.0-alpha.10
$ export IP_ADDRESS=<internal ip-adress>
$ export K8S_VERSION=v1.6.0-alpha.10
$ export USE_CNI=true
$ ./master.sh
```

Shut down with `./turn-down.sh`
