FROM golang:1.13-alpine as builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/hyperjiang/move

COPY . .

RUN go mod download

RUN go build -a -o /tmp/move .

FROM quay.io/hyper/mysql

RUN mkdir /app

RUN mkdir /data

COPY --from=builder /tmp/move /app
COPY ./config.toml /app

ENTRYPOINT ["/app/move"]
