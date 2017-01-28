FROM golang:onbuild

RUN apt-get update && apt-get install -y libpcap-dev

EXPOSE 80
