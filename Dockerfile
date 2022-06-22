FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5
WORKDIR /root
COPY example example/
COPY client.properties .
COPY out/kafka-client kafka-client
RUN chmod +x kafka-client
CMD ["./kafka-client"]