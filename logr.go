package logr

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelDebug
	LevelWarn
	LevelError

	reset = "\033[0m"

	red       = 91
	yellow    = 93
	blueish   = 94
	gray      = 97
	lightGray = 37
	darkGray  = 90
	dimmed    = 2

	labelWidth = 7
)

type Logr struct {
	*log.Logger
	level      LogLevel
	timeFormat string
}

func New(minLevel LogLevel) *Logr {
	return &Logr{
		Logger:     log.New(os.Stdout, "", 0),
		level:      minLevel,
		timeFormat: "",
	}
}

func (l *Logr) WithTime() *Logr {
	return &Logr{
		Logger:     l.Logger,
		level:      l.level,
		timeFormat: "15:04:05",
	}
}

func bold(v string) string {
	return fmt.Sprintf("\033[1m%s", v)
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}
func (l *Logr) log(level LogLevel, colorCode int, label string, msg string) {
	timeStr := ""
	if l.timeFormat != "" {
		timeStr = fmt.Sprintf("%s", time.Now().Format(l.timeFormat))
	}
	if level >= l.level {
		l.Printf("%s %s %s", bold(colorize(colorCode, label)), colorize(darkGray, timeStr), msg)
	}
}

func (l *Logr) Info(msg string) {
	l.log(LevelInfo, blueish, "[INFO]", msg)
}

func (l *Logr) Debug(msg string) {
	l.log(LevelDebug, dimmed, "[DEBUG]", msg)
}
func (l *Logr) Warn(msg string) {
	l.log(LevelWarn, yellow, "[WARN]", msg)
}
func (l *Logr) Error(msg string) {
	l.log(LevelError, red, "[ERROR]", msg)
}

// Used the same as log.Fatalf
//
//	logr.Logf(logr.LevelInfo, "Random: %v", []string{"Hello", "world"})
func (l *Logr) Logf(level LogLevel, format string, v ...any) {
	var colorCode int
	var label string

	switch level {
	case LevelInfo:
		colorCode, label = blueish, "[INFO]"
	case LevelDebug:
		colorCode, label = dimmed, "[DEBUG]"
	case LevelWarn:
		colorCode, label = yellow, "[WARN]"
	case LevelError:
		colorCode, label = red, "[ERROR]"
	}

	timeStr := ""
	if l.timeFormat != "" {
		timeStr = " " + time.Now().Format(l.timeFormat) + " "
	}

	if level >= l.level {
		l.Printf("%s%s  %s", bold(colorize(colorCode, label)),
			colorize(darkGray, timeStr),
			fmt.Sprintf(colorize(lightGray, format), v...))
	}
}
