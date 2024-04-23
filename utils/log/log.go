package log

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var Level LogLevel = INFO

type LogLevel byte

const (
	TRACE LogLevel = 1
	DEBUG LogLevel = 2
	INFO  LogLevel = 3
	WARN  LogLevel = 4
	ERROR LogLevel = 5
	FATAL LogLevel = 6
)

func Trace(format string, args ...interface{}) {
	if TRACE >= Level {
		output(t, format, args...)
	}
}

func Debug(format string, args ...interface{}) {
	if DEBUG >= Level {
		output(d, format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if INFO >= Level {
		output(i, format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if WARN >= Level {
		output(w, format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if ERROR >= Level {
		output(e, format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if FATAL >= Level {
		output(f, format, args...)
	}
}

var (
	t = newLogger("[TRACE] ")
	d = newLogger("[DEBUG] ")
	i = newLogger("[INFOR] ")
	w = newLogger("[WARNI] ")
	e = newLogger("[ERROR] ")
	f = newLogger("[FATAL] ")
)

func newLogger(prefix string) *log.Logger {
	return log.New(os.Stderr, prefix, log.Lshortfile|log.LstdFlags)
}

var lock sync.Mutex

func output(l *log.Logger, format string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	s := fmt.Sprintf(format, args...)
	l.Output(3, s)
}
