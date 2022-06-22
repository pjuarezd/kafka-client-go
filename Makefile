VERSION ?=latest
OWNER ?=pjuarezd
TAG ?= "$(OWNER)/kafka-client:$(VERSION)"
GOTAG ?="$(OWNER)/kafka-client:go"

all: build docker

build:
	@CGO_ENABLED=1 go build --ldflags "-s -w" -o out/kafka-client .

docker:
	@docker buildx build --push -t $(TAG) .
	@docker buildx build --push -t $(GOTAG) .
