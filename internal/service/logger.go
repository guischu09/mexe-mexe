package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
)

type GameLogger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
	level         int
	uuid          string
	module        string
}

func NewLogger(level int, uuid string) *GameLogger {
	flag := 0

	debugOut := io.Discard
	infoOut := io.Discard
	warnOut := io.Discard
	errorOut := io.Discard
	fatalOut := io.Discard

	if level <= LEVEL_DEBUG {
		debugOut = os.Stdout
	}
	if level <= LEVEL_INFO {
		infoOut = os.Stdout
	}
	if level <= LEVEL_WARNING {
		warnOut = os.Stdout
	}
	if level <= LEVEL_ERROR {
		errorOut = os.Stderr
	}

	fatalOut = os.Stderr

	shortUUID := uuid
	if len(uuid) > 4 {
		shortUUID = uuid[:5]
	}

	return &GameLogger{
		debugLogger:   log.New(debugOut, "", flag),
		infoLogger:    log.New(infoOut, "", flag),
		warningLogger: log.New(warnOut, "", flag),
		errorLogger:   log.New(errorOut, "", flag),
		fatalLogger:   log.New(fatalOut, "", flag),
		level:         level,
		uuid:          shortUUID,
		module:        "",
	}
}

func (l *GameLogger) formatLogMessage(message string, calldepth int) string {
	_, file, _, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
	}

	module := filepath.Base(file)
	module = module[:len(module)-len(filepath.Ext(module))]

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	return fmt.Sprintf("%s [%s] %s:: !> %s",
		timestamp,
		l.uuid,
		module,
		message)
}

// Debug logs a message at debug level
func (l *GameLogger) Debug(v ...interface{}) {
	if l.level <= LEVEL_DEBUG {
		message := fmt.Sprint(v...)
		l.debugLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Info logs a message at info level
func (l *GameLogger) Info(v ...interface{}) {
	if l.level <= LEVEL_INFO {
		message := fmt.Sprint(v...)
		l.infoLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Warning logs a message at warning level
func (l *GameLogger) Warning(v ...interface{}) {
	if l.level <= LEVEL_WARNING {
		message := fmt.Sprint(v...)
		l.warningLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Error logs a message at error level
func (l *GameLogger) Error(v ...interface{}) {
	if l.level <= LEVEL_ERROR {
		message := fmt.Sprint(v...)
		l.errorLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Debugf logs a formatted message at debug level
func (l *GameLogger) Debugf(format string, v ...interface{}) {
	if l.level <= LEVEL_DEBUG {
		message := fmt.Sprintf(format, v...)
		l.debugLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Infof logs a formatted message at info level
func (l *GameLogger) Infof(format string, v ...interface{}) {
	if l.level <= LEVEL_INFO {
		message := fmt.Sprintf(format, v...)
		l.infoLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Warningf logs a formatted message at warning level
func (l *GameLogger) Warningf(format string, v ...interface{}) {
	if l.level <= LEVEL_WARNING {
		message := fmt.Sprintf(format, v...)
		l.warningLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Errorf logs a formatted message at error level
func (l *GameLogger) Errorf(format string, v ...interface{}) {
	if l.level <= LEVEL_ERROR {
		message := fmt.Sprintf(format, v...)
		l.errorLogger.Output(3, l.formatLogMessage(message, 2))
	}
}

// Fatal logs a message at fatal level and then terminates the program
func (l *GameLogger) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	l.fatalLogger.Output(3, l.formatLogMessage(message, 2))
	os.Exit(1)
}

// Fatalf logs a formatted message at fatal level and then terminates the program
func (l *GameLogger) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.fatalLogger.Output(3, l.formatLogMessage(message, 2))
	os.Exit(1)
}

// Singleton instance for global access
var GlobalLogger *GameLogger

func SetupLogging(level int, uuid string) {
	GlobalLogger = NewLogger(level, uuid)
}

func GetLogger() *GameLogger {
	if GlobalLogger == nil {
		return NewLogger(LEVEL_INFO, "0000")
	}
	return GlobalLogger
}
