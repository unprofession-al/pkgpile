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
	Command   string `json:"command"`
	Result    string `json:"result"`
}

func (l *Logger) out(cmd string, res string) {
	log := &logmgs{
		Timestamp: time.Now().Format("2006/01/02-15:04:05.000"),
		Command:   cmd,
		Result:    res,
	}

	b, _ := json.Marshal(log)
	l.Println(string(b))
}
