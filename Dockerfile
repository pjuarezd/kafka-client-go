FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5

COPY example .
COPY client.properties .
COPY out/kafka-client kafka-client
