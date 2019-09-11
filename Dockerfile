FROM golang:1.13-alpine as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/github.com/hyperjiang/move

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -a -o /tmp/move .

FROM quay.io/hyper/mysql

RUN mkdir /app
RUN mkdir /data

COPY --from=builder /tmp/move /app
COPY ./config.toml /app

ENTRYPOINT ["/app/move"]
