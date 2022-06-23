VERSION ?=latest
OWNER ?=pjuarezd
TAG ?= "$(OWNER)/kafka-client:$(VERSION)"
GOTAG ?="$(OWNER)/kafka-client:go"

all: deps build docker

deps:
	@go get

build:
	@CGO_ENABLED=1 go build --ldflags "-s -w" -o out/kafka-client .

docker:
	@docker login -u $(OWNER) --password $(DOCKER_HUB_PASSWORD)
	@docker buildx build --push -t $(TAG) .
	@docker buildx build --push -t $(GOTAG) .
