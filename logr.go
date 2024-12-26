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
	fileWidth  = 10
)

type Logr interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
	Logf(level LogLevel, format string, v ...any)
}

type Log struct {
	*log.Logger
	level      LogLevel
	timeFormat string
	showFile   bool
	saveToFile bool
	useColor   bool
}

func New(minLevel ...LogLevel) *Log {
	if len(minLevel) == 0 {
		minLevel = append(minLevel, LevelInfo)
	}
	return &Log{
		Logger:     log.New(os.Stdout, "", 0),
		level:      minLevel[0],
		timeFormat: "",
		saveToFile: false,
		useColor:   false,
	}
}

// Show logs with color
func (l *Log) WithColor() *Log {
	return &Log{
		Logger:     l.Logger,
		level:      l.level,
		timeFormat: l.timeFormat,
		saveToFile: l.saveToFile,
		useColor:   true,
	}
}

// Pass in a time format else defaults to `2006-01-02 15:04:05`
func (l *Log) WithTime(format ...string) *Log {
	if len(format) == 0 {
		format = append(format, "2006-01-02 15:04:05")
	}

	return &Log{
		Logger:     l.Logger,
		level:      l.level,
		timeFormat: format[0],
		saveToFile: l.saveToFile,
		useColor:   l.useColor,
	}
}

func (l *Log) log(level LogLevel, colorCode int, label string, msg string) {
	if level < l.level {
		return
	}

	timeStr := ""
	if l.timeFormat != "" {
		timeStr = time.Now().Format(l.timeFormat)
	}

	paddedLabel := padLabel(label)
	paddedTime := l.padTime(timeStr)

	fileStr := getFileInfo(3)

	if l.useColor {
		l.Printf("%s%s %s%s",
			bold(colorize(colorCode, paddedLabel)),
			colorize(darkGray, fileStr),
			colorize(darkGray, paddedTime),
			msg)
		return
	}

	l.Printf("%s%s %s%s", bold(paddedLabel), fileStr, paddedTime, msg)
}

func (l *Log) Info(msg string) {
	l.log(LevelInfo, blueish, "[INFO]", msg)
}

func (l *Log) Debug(msg string) {
	l.log(LevelDebug, dimmed, "[DEBUG]", msg)
}

func (l *Log) Warn(msg string) {
	l.log(LevelWarn, yellow, "[WARN]", msg)
}

func (l *Log) Error(msg string) {
	l.log(LevelError, red, "[ERROR]", msg)
}

func (l *Log) Logf(level LogLevel, format string, v ...any) {
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

	fileStr := getFileInfo(2)
	if l.useColor {
		l.Printf("%s%s %s%s",
			bold(colorize(colorCode, paddedLabel)),
			colorize(darkGray, fileStr),
			colorize(darkGray, paddedTime),
			fmt.Sprintf(format, v...))

		return
	}

	l.Printf("%s%s %s%s",
		bold(paddedLabel),
		fileStr,
		paddedTime,
		fmt.Sprintf(format, v...),
	)

}
func getFileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return fmt.Sprintf("%-*s", fileWidth, "???:0")
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%-*s", fileWidth, fmt.Sprintf("%s:%d", file, line))
}

func bold(v string) string {
	return fmt.Sprintf("\033[1m%s", v)
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}
func padLabel(label string) string {
	return fmt.Sprintf("%-*s", labelWidth, label)
}

func (l *Log) padTime(timeStr string) string {
	return fmt.Sprintf("%-*s  ", len(l.timeFormat), timeStr)
}

func (l *Log) padFile(file string) string {
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
