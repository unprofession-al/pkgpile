package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "", 0)}
}

type logmgs struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Result    string `json:"result"`
}

func (l *Logger) l(act string, res string) {
	log := &logmgs{
		Timestamp: time.Now().Format("2006/01/02-15:04:05.000"),
		Action:    act,
		Result:    res,
	}

	b, _ := json.Marshal(log)
	l.Println(string(b))
}
