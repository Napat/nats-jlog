# NATS-JLOG(NATS JSON Logger)

NATS-JLOG used to log string to JSON pattern file via NATS protocol.

``` bash
# start nats & nats-jlog server
docker-compose up -d

# basic example client
go run cmd/exclient/exclient.go

# check log
du -h log/*

# logrus example client
go run cmd/exclientLogrus/exclientLogrus.go

# check log
du -h log/*
```

Example log

``` bash
$ cat example.log
{"level":"info","msg":"Hello nats jlogger 2 1234567890","time":"2020-11-08T21:15:25.546+07:00"}
{"level":"info","msg":"Hello nats jlogger 3 1234567890","time":"2020-11-08T21:15:25.548+07:00"}
{"level":"info","msg":"Hello nats jlogger 4 1234567890","time":"2020-11-08T21:15:25.549+07:00"}
{"level":"info","msg":"Hello nats jlogger 5 1234567890","time":"2020-11-08T21:15:25.550+07:00"}
{"level":"info","msg":"Hello nats jlogger 6 1234567890","time":"2020-11-08T21:15:25.551+07:00"}
```

## Task Lists

- [x] flags: natsurl / subject / logPattern / verbose
- [x] auto-sync log file
- [x] dockerfile & compose
- [ ] dockerhub
- [x] example with logrus
- [ ] example with gol
- [ ] NATS TLS

## Reference

- [NATS](https://github.com/nats-io)
- [Logrus](https://github.com/sirupsen/logrus)
