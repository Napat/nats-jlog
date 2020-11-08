package jlogapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// HookNats logrus hooker for Nats client to send log to server
type HookNats struct {
	jl      *JLog
	subject string
	Verbose bool
}

// NewHookNats create a logrus hook for Nats client
func NewHookNats(jl *JLog, subject string) *HookNats {
	hook := HookNats{
		jl:      jl,
		subject: subject,
		Verbose: false,
	}
	return &hook
}

// Fire HookNats to nats loggin server
func (h *HookNats) Fire(entry *logrus.Entry) error {
	s, _ := entry.String()
	if h.Verbose == true {
		fmt.Fprintf(os.Stdout, "%s", s)
	}
	h.natsLogging(s)
	return nil
}

// Levels All logrus level
func (h *HookNats) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *HookNats) natsLogging(jsonStr string) {
	jl := h.jl
	subject := h.subject
	in := &JLogReq{
		Subject: subject,
		Logline: jsonStr,
	}
	inBytes, err := json.Marshal(in)
	if err != nil {
		log.Println(err)
		return
	}
	jl.Pub(in.Subject, inBytes)
}

// VerboseEnable enable/disable print to os.Stdout
func (h *HookNats) VerboseEnable(v bool) {
	h.Verbose = v
}
