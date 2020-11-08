package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Napat/nats-jlog/pkg/jlogapi"
)

func main() {
	fmt.Println("This is client")

	var err error
	natsUrl := "nats://localhost:4222"

	jl, err := jlogapi.JLogNew(jlogapi.MODECLIENT, natsUrl, jlogapi.DefaultLogPattern, 10, "bar", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer jl.Close()

	for {
		jsonStr := `{"time":"` + time.Now().Format("2006-01-02 15:04:05") + `","app":"jlog-client","msg":"looping"}`
		loggingBasic(jl, jsonStr)
		// time.Sleep(1 * time.Second)
		time.Sleep(500 * time.Microsecond)
	}
}

func loggingBasic(jl *jlogapi.JLog, jsonStr string) {
	in := &jlogapi.JLogReq{
		Subject: jlogapi.DefaultSubjectJLog,
		Logline: jsonStr,
	}
	inBytes, err := json.Marshal(in)
	if err != nil {
		log.Println(err)
		return
	}
	jl.Pub(jlogapi.DefaultSubjectJLog, inBytes)
}
