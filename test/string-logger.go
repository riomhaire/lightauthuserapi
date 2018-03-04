package test

import (
	"bytes"
	"fmt"
	"strings"
)

// StringLogger
// This logger is only useful for test where all output is collected
// within a string which can be use to searched for known patterns.
type StringLogger struct {
	log bytes.Buffer
}

func NewStringLogger() StringLogger {
	logger := StringLogger{}

	return logger
}

func (d StringLogger) Log(level, message string) {
	d.log.WriteString(fmt.Sprintf("[%s] %s\n", level, message))
}

func (d StringLogger) Contains(message string) bool {
	return strings.Contains(d.log.String(), message)
}
