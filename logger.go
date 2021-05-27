package log

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"
)

// Create a new Logger with PrintTimeHeader and PrintNilCallerHeader
func NewLogger(w io.Writer) *Logger {
	l := &Logger{Writer: w}
	l.PrintTimeHeader = PrintTimeHeader
	l.PrintCallerHeader = PrintNilCallerHeader
	return l
}

// This logger format is "date_time caller_filepath_line log".
// Example: "2006-01-02 15:04:05.123456 /myproject/awesome.go:127 this is a log."
// You can change time format and caller format by set Logger.PrintTimeHeader and Logger.PrintCallerHeader.
type Logger struct {
	// Where the log output.
	io.Writer
	// Function to format date time header,default is PrintTimeHeader().
	PrintTimeHeader func(*Log)
	// Function to format caller header,default is PrintNilCallerHeader().
	PrintCallerHeader func(*Log, int)
}

// Use fmt.Fprint() to format and add '\n' in the end.
func (l *Logger) Print(a ...interface{}) {
	log := GetLog()
	l.PrintTimeHeader(log)
	l.PrintCallerHeader(log, 2)
	fmt.Fprint(log, a...)
	log.line = append(log.line, '\n')
	l.Writer.Write(log.line)
	PutLog(log)
}

// Use fmt.Fprintf() to format and add '\n' in the end.
func (l *Logger) Printf(f string, a ...interface{}) {
	log := GetLog()
	l.PrintTimeHeader(log)
	l.PrintCallerHeader(log, 2)
	fmt.Fprintf(log, f, a...)
	log.line = append(log.line, '\n')
	l.Writer.Write(log.line)
	PutLog(log)
}

// Format time.Now() into log like "2006-01-02 15:04:05.123456".
func PrintTimeHeader(log *Log) {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// Date
	log.WriteLeftAlignInt(year, 4)
	log.line = append(log.line, '-')
	log.WriteRightAlignInt(int(month), 2)
	log.line = append(log.line, '-')
	log.WriteRightAlignInt(day, 2)
	log.line = append(log.line, ' ')
	// Time
	log.WriteRightAlignInt(hour, 2)
	log.line = append(log.line, ':')
	log.WriteRightAlignInt(minute, 2)
	log.line = append(log.line, ':')
	log.WriteRightAlignInt(second, 2)
	// Nanosecond
	log.line = append(log.line, '.')
	log.WriteLeftAlignInt(t.Nanosecond(), 6)
	// Space
	log.line = append(log.line, ' ')
}

// Do nothing.
func PrintNilCallerHeader(log *Log, skip int) {}

// Format caller file path and line like "/myproject/awesome.go:127" into log.
func PrintFilePathCallerHeader(log *Log, skip int) {
	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		path = "???"
		line = -1
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
	// Space
	log.line = append(log.line, ' ')
}

// Format caller file name and line like "awesome.go:127" into log.
func PrintFileNameCallerHeader(log *Log, skip int) {
	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == filepath.Separator {
				path = path[i+1:]
				break
			}
		}
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
	// Space
	log.line = append(log.line, ' ')
}
