# Kubernetes IoT addon
IoT addon for Kubernetes clusters.

## Status
[![Build Status](https://travis-ci.org/fest-research/iot-addon.svg?branch=master)](https://travis-ci.org/fest-research/iot-addon)

## Quick Start

#### 1. Create insecure Kubernetes cluster

```shell
$ curl https://raw.githubusercontent.com/fest-research/iot-addon/master/assets/hyperkube/hyperkube.sh | sh
```
You can shut down with `docker kill $(docker ps -q)`. Execute twice because some containers might have been restarted by Kubernetes.

Install and configure `kubectl` to connect with master

#### 2. Deploy IoT-Addon

```shell
$ kubectl create -f https://raw.githubusercontent.com/fest-research/iot-addon/master/assets/iot-addon.yaml
```

#### 3. Register RaspberryPIs
Flash RaspberryPi devices with this [software](https://github.com/fest-research/ubikube) to connect easily to the iot-server. 

#### 4. Deploy Demo
Deploy a sample application to the Kubernetes cluster and all RaspberryPis.

```shell
$ kubectl create -f https://raw.githubusercontent.com/fest-research/demo/master/assets/demo-deployment-all.yaml
```
 Cloud part can be found [here](https://github.com/fest-research/demo) and device part [here](https://github.com/fest-research/demo-raspi). Please note: the backend IP is currently hard coded, so you might want to fork the project.


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

