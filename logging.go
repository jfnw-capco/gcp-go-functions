package nozzle

import (
	"encoding/json"
	"fmt"
)

const (
	//DEBUG level
	DEBUG = "DEBUG"

	//INFO level
	INFO = "INFO"

	//WARN level
	WARN = "WARN"

	//ERROR level
	ERROR = "ERROR"
)

var logger = new(simpleLogger)

type simpleLogger struct{}

// LogEntry represents an entry to log with message
type LogEntry struct {
	Level   string            `json:"level"`
	Action  string            `json:"action"`
	Message string            `json:"message,omitempty"`
	Map     map[string]string `json:"data,omitempty"`
	Error   error             `json:"error,omitempty"`
}

// ToString converts an entry to well shaped JSON
func (entry *LogEntry) ToString() string {

	json, err := json.Marshal(entry)
	if err != nil {
		action := fmt.Sprintf("Creating Entry -> %s", entry.Action)
		logger.Error(action, err)
	}

	return string(json)
}

func (l *simpleLogger) Debug(entry LogEntry) {

	entry.Level = DEBUG
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) Info(entry LogEntry) {

	entry.Level = INFO
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) Warn(entry LogEntry) {

	entry.Level = WARN
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) Error(action string, err error) {

	entry := LogEntry{
		Level:  ERROR,
		Action: action,
		Error:  err,
	}

	fmt.Println(entry.ToString())
}
