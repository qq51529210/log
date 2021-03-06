package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	defultLogger = NewLogger(os.Stderr, FormatTimeHeader, FormatFilePathStackHeader)
)

// Create a new Logger.
// If fmtTimeHeader is nil, it will set to formatTimeHeader.
// If fmtStackHeader is nil, it will set to PrintNilCallerHeader.
func NewLogger(out io.Writer, fmtTimeHeader func(*Log), fmtStackHeader func(*Log, int)) *Logger {
	l := new(Logger)
	l.out = out
	l.SetFormatTimeHeader(fmtTimeHeader)
	l.SetFormatStackHeader(fmtStackHeader)
	return l
}

// Format is "2006-01-02 15:04:05.123456".
func FormatTimeHeader(log *Log) {
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
}

// Format is "/myproject/awesome.go:127".
func FormatFilePathStackHeader(log *Log, skip int) {
	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		path = "???"
		line = -1
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
}

// Format is "awesome.go:127".
func FormatFileNameStackHeader(log *Log, skip int) {
	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if path[i] == filepath.Separator {
				path = path[i+1:]
				break
			}
		}
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
}

type Logger struct {
	// Output log.
	out io.Writer
	// Format date time header.
	fmtTimeHeader func(*Log)
	// Format code stack header.
	fmtStackHeader func(*Log, int)
}

func (l *Logger) SetFormatTimeHeader(f func(*Log)) {
	if f != nil {
		l.fmtTimeHeader = f
	} else {
		l.fmtTimeHeader = func(*Log) {}
	}
}

func (l *Logger) SetFormatStackHeader(f func(*Log, int)) {
	if f != nil {
		l.fmtStackHeader = f
	} else {
		l.fmtStackHeader = func(*Log, int) {}
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.out = w
}

// Format is "fmtTimeHeader() fmtStackHeader() s\n".
func (l *Logger) Print(s string) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, 2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	log.line = append(log.line, s...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "fmtTimeHeader() fmtStackHeader() fmt.Fprintln()".
func (l *Logger) Fprint(a ...interface{}) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, 2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintln(log, a...)
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "fmtTimeHeader() fmtStackHeader() fmt.Fprintf()\n".
func (l *Logger) Fprintf(f string, a ...interface{}) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, 2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintf(log, f, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "fmtTimeHeader() fmtStackHeader(skip) s\n".
func (l *Logger) PrintStack(skip int, s string) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, skip+2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	log.line = append(log.line, s...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "fmtTimeHeader() fmtStackHeader(skip) fmt.Fprintln()".
func (l *Logger) FprintStack(skip int, a ...interface{}) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, skip+2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintln(log, a...)
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "fmtTimeHeader() fmtStackHeader(skip) fmt.Fprintf()\n".
func (l *Logger) FprintfStack(skip int, f string, a ...interface{}) {
	log := GetLog()
	// Time
	l.fmtTimeHeader(log)
	// Space
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, skip+2)
	// Space
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintf(log, f, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "v fmtTimeHeader() fmtStackHeader(skip) fmt.Fprintln()".
func (l *Logger) printLevel(v byte, skip int, a ...interface{}) {
	log := GetLog()
	// Level
	log.line = append(log.line, v)
	log.line = append(log.line, ' ')
	// Time
	l.fmtTimeHeader(log)
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, skip)
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintln(log, a...)
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

// Format is "v fmtTimeHeader() fmtStackHeader(skip) fmt.Fprintf()\n".
func (l *Logger) printfLevel(v byte, skip int, f string, a ...interface{}) {
	log := GetLog()
	// Level
	log.line = append(log.line, v)
	log.line = append(log.line, ' ')
	// Time
	l.fmtTimeHeader(log)
	log.line = append(log.line, ' ')
	// Caller
	l.fmtStackHeader(log, skip)
	log.line = append(log.line, ' ')
	// Log
	fmt.Fprintf(log, f, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.out.Write(log.line)
	PutLog(log)
}

func (l *Logger) Debug(a ...interface{}) {
	l.printLevel('D', 3, a...)
}

func (l *Logger) DebugStack(skip int, a ...interface{}) {
	l.printLevel('D', skip+3, a...)
}

func (l *Logger) Debugf(f string, a ...interface{}) {
	l.printfLevel('D', 3, f, a...)
}

func (l *Logger) DebugfStack(skip int, f string, a ...interface{}) {
	l.printfLevel('D', skip+3, f, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.printLevel('I', 3, a...)
}

func (l *Logger) InfoStack(skip int, a ...interface{}) {
	l.printLevel('I', skip+3, a...)
}

func (l *Logger) Infof(f string, a ...interface{}) {
	l.printfLevel('I', 3, f, a...)
}

func (l *Logger) InfofStack(skip int, f string, a ...interface{}) {
	l.printfLevel('I', skip+3, f, a...)
}

func (l *Logger) Warn(a ...interface{}) {
	l.printLevel('W', 3, a...)
}

func (l *Logger) WarnStack(skip int, a ...interface{}) {
	l.printLevel('W', skip+3, a...)
}

func (l *Logger) Warnf(f string, a ...interface{}) {
	l.printfLevel('W', 3, f, a...)
}

func (l *Logger) WarnfStack(skip int, f string, a ...interface{}) {
	l.printfLevel('W', skip+3, f, a...)
}

func (l *Logger) Error(a ...interface{}) {
	l.printLevel('E', 3, a...)
}

func (l *Logger) ErrorStack(skip int, a ...interface{}) {
	l.printLevel('E', skip+3, a...)
}

func (l *Logger) Errorf(f string, a ...interface{}) {
	l.printfLevel('E', 3, f, a...)
}

func (l *Logger) ErrorfStack(skip int, f string, a ...interface{}) {
	l.printfLevel('E', skip+3, f, a...)
}

func SetFormatStackHeader(formatStackHeader func(*Log, int)) {
	defultLogger.SetFormatStackHeader(formatStackHeader)
}

func SetFormatTimeHeader(formatTimeHeader func(*Log)) {
	defultLogger.SetFormatTimeHeader(formatTimeHeader)
}

func SetOutput(out io.Writer) {
	defultLogger.out = out
}

func Print(s string) {
	defultLogger.PrintStack(1, s)
}

func Fprint(a ...interface{}) {
	defultLogger.FprintStack(1, a...)
}

func Fprintf(f string, a ...interface{}) {
	defultLogger.FprintfStack(1, f, a...)
}

func Debug(a ...interface{}) {
	defultLogger.printLevel('D', 3, a...)
}

func DebugStack(skip int, a ...interface{}) {
	defultLogger.printLevel('D', skip+3, a...)
}

func Debugf(f string, a ...interface{}) {
	defultLogger.printfLevel('D', 3, f, a...)
}

func DebugfStack(skip int, f string, a ...interface{}) {
	defultLogger.printfLevel('D', skip+3, f, a...)
}

func Info(a ...interface{}) {
	defultLogger.printLevel('I', 3, a...)
}

func InfoStack(skip int, a ...interface{}) {
	defultLogger.printLevel('I', skip+3, a...)
}

func Infof(f string, a ...interface{}) {
	defultLogger.printfLevel('I', 3, f, a...)
}

func InfofStack(skip int, f string, a ...interface{}) {
	defultLogger.printfLevel('I', skip+3, f, a...)
}

func Warn(a ...interface{}) {
	defultLogger.printLevel('W', 3, a...)
}

func WarnStack(skip int, a ...interface{}) {
	defultLogger.printLevel('W', skip+3, a...)
}

func Warnf(f string, a ...interface{}) {
	defultLogger.printfLevel('W', 3, f, a...)
}

func WarnfStack(skip int, f string, a ...interface{}) {
	defultLogger.printfLevel('W', skip+3, f, a...)
}

func Error(a ...interface{}) {
	defultLogger.printLevel('E', 3, a...)
}

func ErrorStack(skip int, a ...interface{}) {
	defultLogger.printLevel('E', skip+3, a...)
}

func Errorf(f string, a ...interface{}) {
	defultLogger.printfLevel('E', 3, f, a...)
}

func ErrorfStack(skip int, f string, a ...interface{}) {
	defultLogger.printfLevel('E', skip+3, f, a...)
}
