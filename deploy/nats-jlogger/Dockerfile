FROM golang:1.15-alpine3.12 AS builder
    ENV GOOS=linux
    ENV GOARCH=amd64
    ENV CGO_ENABLED=0
    RUN apk --no-cache add git
    COPY . /go/src/app
    WORKDIR /go/src/app
    RUN go build -a -installsuffix cgo -o nats-jlogger -ldflags "-w -s" cmd/nats-jlogger/nats-jlogger.go

FROM alpine:3.12 AS release
    RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    echo "Asia/Bangkok" >  /etc/timezone
    RUN apk --no-cache add ca-certificates
    WORKDIR /app/
    COPY --from=builder /go/src/app/nats-jlogger ./nats-jlogger
    ENTRYPOINT ["/app/nats-jlogger"]
    CMD [""]

FROM release