package logger

import (
	"fmt"
	"strings"
	"sync"
)

// UILogger is a logger specifically designed to hold logs until they are ready to be dumped all at once.
type UILogger struct {
	logs []string
	lock sync.RWMutex
}

// NewUILogger creates a UILogger instance.
func NewUILogger() *UILogger {
	return &UILogger{}
}

// Write a new log message.  The component should be an identifier for the UI component generating the log.
func (l *UILogger) Write(component, msg string, args ...interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logs = append(l.logs, component+": "+fmt.Sprintf(msg, args...))
}

// Flush all logs in the log buffer.
func (l *UILogger) Flush() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logs = l.logs[0:0]
}

// Dump all of the logs.  Log lines will be separated by \n.
func (l *UILogger) Dump() string {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return strings.Join(l.logs, "\n")
}
