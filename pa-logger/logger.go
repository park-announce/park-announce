package logger

import (
	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type Logger struct {
	*logrus.Logger
	fields logrus.Fields
}

// NewLogger initializes the standard logger
func NewLogger(appname string) *Logger {
	var baseLogger = logrus.New()
	standardFields := logrus.Fields{
		"appname": appname,
	}
	var standardLogger = &Logger{baseLogger, standardFields}

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
	generalMessage         = Event{4, "%s"}
)

// InvalidArg is a standard error message
func (l *Logger) InvalidArg(argumentName string) {
	l.WithFields(l.fields).Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func (l *Logger) InvalidArgValue(argumentName string, argumentValue string) {
	l.WithFields(l.fields).Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (l *Logger) MissingArg(argumentName string) {
	l.WithFields(l.fields).Errorf(missingArgMessage.message, argumentName)
}

func (l *Logger) LogError(message string) {
	l.WithFields(l.fields).Errorf(generalMessage.message, message)
}

func (l *Logger) LogWarning(message string) {
	l.WithFields(l.fields).Warningf(generalMessage.message, message)
}

func (l *Logger) LogInfo(message string) {
	l.WithFields(l.fields).Infof(generalMessage.message, message)
}
