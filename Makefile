VERSION ?=latest
TAG ?= "kafka-client:$(VERSION)"

all: build docker

build:
	@CGO_ENABLED=1 go build --ldflags "-s -w" -o out/kafka-client .

docker:
	@docker build --no-cache -t $(TAG) .
