FROM golang:1.15

WORKDIR /app

ADD . .

RUN go mod download && \
    go build

CMD ["./example-go"]
