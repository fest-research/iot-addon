GO := $(shell command -v go 2> /dev/null)

prepare:
	@mkdir -p ./build/apiserver
	@mkdir -p ./build/controller

all: prepare apiserver controller

apiserver: prepare
ifndef GO
	$(error "Could not find GO compiler.")
endif
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/apiserver/apiserver cmd/apiserver/apiserver.go

controller: prepare
ifndef GO
	$(error "Could not find GO compiler.")
endif
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/controller/controller cmd/controller/controller.go

clean:
	@rm -rf ./build/apiserver/apiserver
	@rm -rf ./build/controller/controller