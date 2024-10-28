package logr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota + 1
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

	labelWidth = 8
	fileWidth  = 20

	skip = 3
)

type Logr struct {
	*log.Logger
	level      LogLevel
	timeFormat string
	showFile   bool
	saveToFile bool
}

func New(minLevel LogLevel) *Logr {
	return &Logr{
		Logger:     log.New(os.Stdout, "", 0),
		level:      minLevel,
		timeFormat: "",
		saveToFile: false,
	}
}

// Pass in a time format else defaults to `2006-01-02 15:04:05`
func (l *Logr) WithTime(format ...string) *Logr {
	if len(format) == 0 {
		format = append(format, "2006-01-02 15:04:05")
	}

	return &Logr{
		Logger:     l.Logger,
		level:      l.level,
		timeFormat: format[0],
		saveToFile: l.saveToFile,
	}
}

func (l *Logr) WithFile() *Logr {
	return &Logr{
		Logger:     l.Logger,
		level:      l.level,
		timeFormat: l.timeFormat,
		saveToFile: false,
	}
}

func (l *Logr) log(level LogLevel, colorCode int, label string, msg string) {
	if level < l.level {
		return
	}

	timeStr := ""
	if l.timeFormat != "" {
		timeStr = time.Now().Format(l.timeFormat)
	}

	paddedLabel := padLabel(label)
	paddedTime := l.padTime(timeStr)

	fileStr := getFileInfo()

	l.Printf("%s%s %s%s",
		bold(colorize(colorCode, paddedLabel)),
		colorize(darkGray, fileStr),
		colorize(darkGray, paddedTime),
		msg)
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

func (l *Logr) Logf(level LogLevel, format string, v ...any) {
	if level < l.level {
		return
	}

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
	default:
		colorCode, label = blueish, "[INFO]"
	}

	timeStr := ""
	if l.timeFormat != "" {
		timeStr = time.Now().Format(l.timeFormat)
	}

	paddedLabel := padLabel(label)
	paddedTime := l.padTime(timeStr)

	l.Printf("%s%s%s",
		bold(colorize(colorCode, paddedLabel)),
		colorize(darkGray, paddedTime),
		fmt.Sprintf(colorize(lightGray, format), v...))
}

func getFileInfo() string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???:0"
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}

func bold(v string) string {
	return fmt.Sprintf("\033[1m%s", v)
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func padLabel(label string) string {
	visibleLen := len(label)
	padding := strings.Repeat(" ", (len(label)+1)-visibleLen)
	return label + padding
}

func (l *Logr) padTime(timeStr string) string {
	if timeStr == "" {
		return ""
	}

	padding := strings.Repeat(" ", (len(l.timeFormat)+2)-len(timeStr))
	return timeStr + padding
}

func (l *Logr) padFile(file string) string {
	if file == "" {
		return ""
	}
	// -3 to account for "..."
	if len(file) > fileWidth-3 {
		file = "..." + file[len(file)-(fileWidth-3):]
	}
	padding := strings.Repeat(" ", fileWidth-len(file))
	return file + padding
}
