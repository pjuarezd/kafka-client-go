FROM ubuntu:latest
WORKDIR /root
COPY example example/
COPY client.properties .
COPY out/kafka-client kafka-client

RUN \
    apt-get update && \
    apt-get install -y curl ca-certificates golang-go

CMD ["./kafka-client"]