GO := $(shell command -v go 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)

DOCKER_HUB = fest

prepare:
	@mkdir -p ./build/apiserver
	@mkdir -p ./build/controller

check_docker:
ifndef DOCKER
	$(error "Could not find docker.")
endif

check_go:
ifndef GO
	$(error "Could not find GO compiler.")
endif

build: check_go prepare apiserver controller

apiserver: prepare

	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/apiserver/apiserver cmd/apiserver/apiserver.go

controller: prepare
ifndef GO
	$(error "Could not find GO compiler.")
endif
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/controller/controller cmd/controller/controller.go

clean:
	@rm -rf ./build/apiserver/apiserver
	@rm -rf ./build/controller/controller

build_docker: check_go
	docker build -t $(DOCKER_HUB)/iot-apiserver build/apiserver
	docker build -t $(DOCKER_HUB)/iot-controller build/controller

deploy: build_docker
	docker push $(DOCKER_HUB)/iot-apiserver
	docker push $(DOCKER_HUB)/iot-controller