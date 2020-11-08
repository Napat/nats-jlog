package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Napat/nats-jlog/pkg/jlogapi"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("This is client")

	var err error
	natsUrl := "nats://localhost:4222"

	// init jlog
	jl, err := jlogapi.JLogNew(jlogapi.MODECLIENT, natsUrl, jlogapi.DefaultLogPattern, 10, "bar", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer jl.Close()

	hook := jlogapi.NewHookNats(jl, jlogapi.DefaultSubjectJLog)
	hook.VerboseEnable(false)

	// init logrus with json format
	var log = logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // add timestamp with format: RFC3339 Micro
		PrettyPrint:     false,
	}
	log.SetOutput(ioutil.Discard) // disable logrus printer
	log.AddHook(hook)

	// logging looper
	cnt := 0
	for {
		cnt++
		log.Infof("Hello nats jlogger %v 1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", cnt)
		time.Sleep(1000 * time.Microsecond)
	}
}
