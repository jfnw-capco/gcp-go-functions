package gcpfunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

var logger = new(simpleLogger)

type simpleLogger struct{}

// LogDataEntry represents an entry to log with data
type LogDataEntry struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data,omitempty"`
	Err    error             `json:"err,omitempty"`
}

// LogMessageEntry represents an entry to log with message
type LogMessageEntry struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

// ToString converts an entry to wll shaped JSON
func (entry *LogDataEntry) ToString() string {

	json, err := json.Marshal(entry)
	if err != nil {
		action := fmt.Sprintf("Creating Entry -> %s", entry.Action)
		logger.Fatal(action, err)
	}

	return string(json)
}

// ToString converts an entry to wll shaped JSON
func (entry *LogMessageEntry) ToString() string {

	output, err := json.Marshal(entry)
	if err != nil {
		action := fmt.Sprintf("Creating Entry -> %s", entry.Action)
		logger.Fatal(action, err)
	}

	compactedBuffer := new(bytes.Buffer)
	err = json.Compact(compactedBuffer, output)
	if err != nil {
		action := fmt.Sprintf("Compacting Entry -> %s", entry.Action)
		logger.Fatal(action, err)
	}

	return string(output)
}

func (l *simpleLogger) DebugMessage(action string, message string) {

	entry := getMessageEntry(action, message, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) InfoMessage(action string, message string) {

	entry := getMessageEntry(action, message, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) WarnMessage(action string, message string) {

	entry := getMessageEntry(action, message, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) DebugData(action string, data map[string]string) {

	entry := getDataEntry(action, data, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) InfoData(action string, data map[string]string) {

	entry := getDataEntry(action, data, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) WarnData(action string, data map[string]string) {

	entry := getDataEntry(action, data, nil)
	fmt.Println(entry.ToString())
}

func (l *simpleLogger) Fatal(action string, err error) {

	log.Fatalln(err)
}

func getMessageEntry(action string, message string, err error) LogMessageEntry {

	return LogMessageEntry{
		Action:  action,
		Message: message,
		Err:     err,
	}
}

func getDataEntry(action string, data map[string]string, err error) LogDataEntry {

	return LogDataEntry{
		Action: action,
		Data:   data,
		Err:    err,
	}
}
