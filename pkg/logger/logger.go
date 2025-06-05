package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
	Debug(msg string)
	WithField(key, value string) Logger
}

type logger struct {
	level  string
	fields map[string]string
}

func New(level string) Logger {
	return &logger{
		level:  level,
		fields: make(map[string]string),
	}
}

func (l *logger) Info(msg string) {
	l.log("INFO", msg)
}

func (l *logger) Error(msg string) {
	l.log("ERROR", msg)
}

func (l *logger) Fatal(msg string) {
	l.log("FATAL", msg)
	os.Exit(1)
}

func (l *logger) Debug(msg string) {
	if l.level == "debug" {
		l.log("DEBUG", msg)
	}
}

func (l *logger) WithField(key, value string) Logger {
	newFields := make(map[string]string)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	
	return &logger{
		level:  l.level,
		fields: newFields,
	}
}

func (l *logger) log(level, msg string) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000Z")
	fieldsStr := ""
	for k, v := range l.fields {
		fieldsStr += fmt.Sprintf(" %s=%s", k, v)
	}
	log.Printf("[%s] %s %s%s", level, timestamp, msg, fieldsStr)
}
