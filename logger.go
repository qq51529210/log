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
	defultLogger = NewLogger(os.Stderr)
)

// Create a new Logger.
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
	// Function to format date time header, default is PrintTimeHeader.
	PrintTimeHeader func(*Log)
	// Function to format caller header, default is PrintNilCallerHeader.
	PrintCallerHeader func(*Log, int)
}

// Print format is "time caller log\n", use fmt.Fprint to format a.
func (l *Logger) print(n int, a ...interface{}) {
	log := GetLog()
	// Time
	l.PrintTimeHeader(log)
	// Caller
	l.PrintCallerHeader(log, n)
	// Log
	fmt.Fprint(log, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.Writer.Write(log.line)
	PutLog(log)
}

// Print format is "time caller log\n", use fmt.Fprintf to format f and a.
func (l *Logger) printf(n int, f string, a ...interface{}) {
	log := GetLog()
	// Time
	l.PrintTimeHeader(log)
	// Caller
	l.PrintCallerHeader(log, n)
	// Log
	fmt.Fprintf(log, f, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.Writer.Write(log.line)
	PutLog(log)
}

// Print format is "level time caller log\n", use fmt.Fprint to format a.
func (l *Logger) printLevel(n int, c byte, a ...interface{}) {
	log := GetLog()
	// Level
	log.line = append(log.line, c)
	log.line = append(log.line, ' ')
	// Time
	l.PrintTimeHeader(log)
	// Caller
	l.PrintCallerHeader(log, n)
	// Log
	fmt.Fprint(log, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.Writer.Write(log.line)
	PutLog(log)
}

// Print format is "level time caller log\n", use fmt.Fprintf to format f and a.
func (l *Logger) printfLevel(n int, c byte, f string, a ...interface{}) {
	log := GetLog()
	// Level
	log.line = append(log.line, c)
	log.line = append(log.line, ' ')
	// Time
	l.PrintTimeHeader(log)
	// Caller
	l.PrintCallerHeader(log, n)
	// Log
	fmt.Fprintf(log, f, a...)
	// Endline
	log.line = append(log.line, '\n')
	// Output
	l.Writer.Write(log.line)
	PutLog(log)
}

// Print format is "time caller log\n".
func (l *Logger) Print(a ...interface{}) {
	l.print(3, a...)
}

// Print format is "time caller log\n".
func (l *Logger) Printf(f string, a ...interface{}) {
	l.printf(3, f, a...)
}

// Print format is "D time caller log\n".
func (l *Logger) Debug(a ...interface{}) {
	l.printLevel(3, 'D', a...)
}

// Print format is "D time caller log\n".
func (l *Logger) Debugf(f string, a ...interface{}) {
	l.printfLevel(3, 'D', f, a...)
}

// Print format is "I time caller log\n".
func (l *Logger) Info(a ...interface{}) {
	l.printLevel(3, 'I', a...)
}

// Print format is "I time caller log\n".
func (l *Logger) Infof(f string, a ...interface{}) {
	l.printfLevel(3, 'I', f, a...)
}

// Print format is "W time caller log\n".
func (l *Logger) Warn(a ...interface{}) {
	l.printLevel(3, 'W', a...)
}

// Print format is "W time caller log\n".
func (l *Logger) Warnf(f string, a ...interface{}) {
	l.printfLevel(3, 'W', f, a...)
}

// Print format is "E time caller log\n".
func (l *Logger) Error(a ...interface{}) {
	l.printLevel(3, 'E', a...)
}

// Print format is "E time caller log\n".
func (l *Logger) Errorf(f string, a ...interface{}) {
	l.printfLevel(3, 'E', f, a...)
}

// Print format is "2006-01-02 15:04:05.123456".
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

// Print format is "/myproject/awesome.go:127".
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

// Print format is "awesome.go:127".
func PrintFileNameCallerHeader(log *Log, skip int) {
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
	// Space
	log.line = append(log.line, ' ')
}

// Set defaultLogger.PrintCallerHeader.
func SetPrintCallerHeader(printCallerHeader func(*Log, int)) {
	defultLogger.PrintCallerHeader = printCallerHeader
}

// Set defaultLogger.PrintTimeHeader.
func SetPrintTimeHeader(printTimeHeader func(*Log)) {
	defultLogger.PrintTimeHeader = printTimeHeader
}

// Set ddefaultLogger.Writer.
func SetWriter(w io.Writer) {
	defultLogger.Writer = w
}

// defaultLogger.Print
func Print(a ...interface{}) {
	defultLogger.print(3, a...)
}

// defaultLogger.Printf
func Printf(f string, a ...interface{}) {
	defultLogger.printf(3, f, a...)
}

// defaultLogger.Debug
func Debug(a ...interface{}) {
	defultLogger.printLevel(3, 'D', a...)
}

// defaultLogger.Debugf
func Debugf(f string, a ...interface{}) {
	defultLogger.printfLevel(3, 'D', f, a...)
}

// defaultLogger.Info
func Info(a ...interface{}) {
	defultLogger.printLevel(3, 'I', a...)
}

// defaultLogger.Infof
func Infof(f string, a ...interface{}) {
	defultLogger.printfLevel(3, 'I', f, a...)
}

// defaultLogger.Warn
func Warn(a ...interface{}) {
	defultLogger.printLevel(3, 'W', a...)
}

// defaultLogger.Warnf
func Warnf(f string, a ...interface{}) {
	defultLogger.printfLevel(3, 'W', f, a...)
}

// defaultLogger.Error
func Error(a ...interface{}) {
	defultLogger.printLevel(3, 'E', a...)
}

// defaultLogger.Warnf
func Errorf(f string, a ...interface{}) {
	defultLogger.printfLevel(3, 'E', f, a...)
}
