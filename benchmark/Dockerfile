FROM ubuntu:latest

MAINTAINER Imran Pochi <imran@kinvolk.io>

RUN apt-get update && \
    apt-get install -y iperf3 ipset iputils-ping

COPY benchmark /usr/bin/benchmark

ENTRYPOINT ["/usr/bin/benchmark"]

