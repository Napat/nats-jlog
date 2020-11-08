package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Napat/nats-jlog/pkg/jlogapi"

	nats "github.com/nats-io/nats.go"
)

var jl *jlogapi.JLog
var logDirty bool
var syncFileTimeSec time.Duration = 1
var verbose bool

var syncMU sync.Mutex

func main() {
	var err error
	var natsUrl string
	var subject string
	var logPattern string

	finished := make(chan bool)

	flag.StringVar(&natsUrl, "natsurl", "nats://localhost:4222", "NATS URL.")
	flag.StringVar(&subject, "subject", "jlog.default", "NATS Subject to log.")
	flag.StringVar(&logPattern, "logfile", "./log/%Y%m%d_default.log", "NATS log pattern.")
	flag.BoolVar(&verbose, "verbose", false, "verbose log output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of [%s]:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	jl, err = jlogapi.JLogNew(jlogapi.MODESERVER, natsUrl, logPattern, 10, "foo", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer jl.Close()

	go logSync(jl)

	jl.QSub(subject, "queue"+subject, cbLogging) //jl.Sub(subject, cbLogging)

	// stamp log starting
	jsonStrIN := `{"time":"` + time.Now().Format("2006-01-02 15:04:05") + `","app":"jlog-server","msg":"starting"}`
	in := &jlogapi.JLogReq{
		Subject: subject,
		Logline: jsonStrIN,
	}
	inBytes, err := json.Marshal(in)
	if err != nil {
		log.Println(err)
		return
	}
	jl.Pub(subject, inBytes)

	<-finished
}

func logSync(jl *jlogapi.JLog) {
	for {
		if logDirty == true {
			if verbose == true {
				log.Println("Sync() log file")
			}

			syncMU.Lock()
			jl.Logwritter.Sync()
			logDirty = false
			syncMU.Unlock()
		}
		time.Sleep(syncFileTimeSec * time.Second)
	}
}

func cbLogging(m *nats.Msg) {
	var jlog jlogapi.JLogReq

	// fmt.Printf("Received a message: %s\n", string(m.Data))
	if err := json.Unmarshal(m.Data, &jlog); err != nil {
		log.Println(err)
		return
	}

	// write log file
	syncMU.Lock()
	if jlog.Logline[len(jlog.Logline)-1] == '\n' {
		if verbose == true {
			log.Print(jlog.Logline)
		}
		jl.Logwritter.Write([]byte(jlog.Logline))
	} else {
		if verbose == true {
			log.Println(jlog.Logline)
		}
		jl.Logwritter.Write([]byte(jlog.Logline + "\n"))
	}
	logDirty = true
	syncMU.Unlock()
}
