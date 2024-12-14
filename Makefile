DOCKER_IMAGE_TAG := $(shell git describe --abbrev=0 --tags)
DOCKER_REGISTRY ?= docker.io
DOCKER_IMAGE_NAME ?= gsapi

.PHONY: build_local
build_local: 
	go mod tidy
	go build -o ./cmd/bin/gsapi ./cmd/main.go

.PHONY: clean_local
clean_local:	
	rm -rf ./cmd/bin


.PHONY: docker_build 
docker_build:
	go mod tidy
	docker build -t ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:latest .


.PHONY: test
test:
	go test -race -v ./...