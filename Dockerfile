FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

RUN apt-get update && apt-get install -y libpcap-dev

EXPOSE 80

CMD ["/app/main"]
