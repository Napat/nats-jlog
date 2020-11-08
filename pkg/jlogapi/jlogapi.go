package jlogapi

import (
	"fmt"
	"log"
	"time"

	"github.com/Napat/go-cronowriter" //"github.com/utahta/go-cronowriter"
	nats "github.com/nats-io/nats.go"
)

const DefaultUrl = "nats://localhost:4222"
const DefaultSubjectJLog = "jlog.default"
const DefaultLogPattern = "./log/%Y%m%d_default.log"

type JLogMode string

const (
	MODECLIENT JLogMode = "JLOGCLIENT"
	MODESERVER          = "JLOGSERVER"
)

type JLogReq struct {
	Subject string `json:"subject"`
	LogPath string `json:"path"`
	Logline string `json:"line"`
}

type JLog struct {
	NC         *nats.Conn
	Logwritter *cronowriter.CronoWriter
	logPattern string
}

// JLogNew create NATS server connection
//	mode: MODESERVER initial connection with Logwritter, MODECLIENT without Logwritter
//  url: nats server. for example,
//		"nats.DefaultURL", "nats://demo.nats.io:4222", "demo.nats.io:4222"
//		or with plaintext user/password "myname:password@127.0.0.1"
//		or with token "mytoken@localhost"
//	logPattern: log pattern, ie., ./log/%Y%m%d_default.log
//  timeoutSec: time limits how long it can take to establish a connection to a server.(e.g., 10)
//  friendlyName: (optional) but is highly recommended as a friendly connection name will help in monitoring, error reporting, debugging, and testing.
//  disconnectErrCB: ** NATS is enable auto-reconnect** use to be notified of disconnect events
//  reconnectCB: ** NATS is enable auto-reconnect** use to be notified of reconnect events
func JLogNew(
	mode JLogMode,
	url string,
	logPattern string,
	timeoutSec int,
	friendlyName string,
	disconnectErrCB func(nc *nats.Conn, err error),
	reconnectCB func(nc *nats.Conn)) (*JLog, error) {

	var jl JLog

	if url == "" {
		return nil, fmt.Errorf("Invalid NATS url")
	}

	opts := nats.GetDefaultOptions()
	opts.Url = url
	opts.Timeout = time.Duration(timeoutSec) * time.Second

	if disconnectErrCB != nil {
		opts.DisconnectedErrCB = disconnectErrCB
	}
	if reconnectCB != nil {
		opts.ReconnectedCB = reconnectCB
	}

	if friendlyName != "" {
		opts.Name = friendlyName
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	// defer nc.Close()

	//log.Printf("Connected to  with status: %v\n", url, nc.Status())
	if nc.Status() != nats.CONNECTED {
		return nil, fmt.Errorf("nats connection to %v error with status: %v\n", url, nc.Status())
	}

	// https://docs.nats.io/developing-with-nats/connecting/misc
	mp := nc.MaxPayload()
	log.Printf("NATS server support maximum payload %v bytes", mp)

	jl.NC = nc

	if mode == MODESERVER {
		jl.Logwritter = cronowriter.MustNew(logPattern) // jl.Logwritter = cronowriter.MustNew("./log/example_%Y%m%d.log")
	}

	return &jl, nil
}

func (jl *JLog) Close() {
	jl.NC.Close()
}

func (jl *JLog) Sub(subject string, cb func(m *nats.Msg)) error {
	// Simple Async Subscribe
	if _, err := jl.NC.Subscribe(subject, cb); err != nil {
		log.Println("Failed to parse protobuf message:", err)
		return err
	}

	return nil
}

func (jl *JLog) QSub(subject string, queue string, cb func(m *nats.Msg)) error {
	// Async Queue Subscribe
	if _, err := jl.NC.QueueSubscribe(subject, queue, cb); err != nil {
		log.Println("Failed to parse protobuf message:", err)
		return err
	}

	return nil
}

func (jl *JLog) Pub(subject string, data []byte) {
	err := jl.NC.Publish(subject, data)
	if err != nil {
		log.Fatal(err)
	}
}
