package log

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
)

type Level byte

const (
	DEBUG_LEVEL Level = 'D'
	INFO_LEVEL  Level = 'I'
	WARN_LEVEL  Level = 'W'
	ERROR_LEVEL Level = 'E'
	PANIC_LEVEL Level = 'P'
)

const (
	_DEFAULT_DEPTH = 3
	_LOGGER_DEPTH  = 3
)

type Logger interface {
	// Set header formater
	SetHeader(header Header)
	// Set output writer
	SetWriter(writer io.Writer)
	// Set output types.
	// For example: SetLevel(ERROR, ERROR, PANIC),
	// DebugXX() and InfoXX() will not output.
	SetLevel(levels ...Level)
	// AppId D Header log
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DepthDebug(depth int, args ...interface{})
	DepthDebugf(depth int, format string, args ...interface{})
	// AppId I Header log
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	DepthInfo(depth int, args ...interface{})
	DepthInfof(depth int, format string, args ...interface{})
	// AppId W Header log
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	DepthWarn(depth int, args ...interface{})
	DepthWarnf(depth int, format string, args ...interface{})
	// AppId E Header log
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	DepthError(depth int, args ...interface{})
	DepthErrorf(depth int, format string, args ...interface{})
	// AppId P Header log
	Recover(recover interface{})
}

var (
	defaultLogger *logger
	runtimePanic  = []byte("runtime/panic.go")
)

func init() {
	// init node with MAC address
	addr, err := net.Interfaces()
	if nil != err {
		panic(err)
	}
	id := ""
	for i := 0; i < len(addr); i++ {
		if addr[i].Flags|net.FlagUp != 0 && len(addr[i].HardwareAddr) != 0 {
			id = addr[i].HardwareAddr.String()
			break
		}
	}
	defaultLogger = &logger{
		Writer: os.Stdout,
		Header: &FilePathStackHeader{
			DefaultHeader: DefaultHeader{
				id: id,
			},
		},
		enableDebug: true,
		enableInfo:  true,
		enableWarn:  true,
		enableError: true,
	}
}

func NewLogger(writer io.Writer, level int, header Header, levels ...Level) Logger {
	lg := new(logger)
	lg.Writer = writer
	lg.Header = header
	lg.SetLevel(levels...)
	return lg
}

func SetWriter(writer io.Writer) {
	defaultLogger.SetWriter(writer)
}

func SetHeader(header Header) {
	defaultLogger.SetHeader(header)
}

func SetLevel(levels ...Level) {
	defaultLogger.SetLevel(levels...)
}

func Debug(args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(_DEFAULT_DEPTH, DEBUG_LEVEL, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(_DEFAULT_DEPTH, DEBUG_LEVEL, format, args...)
	}
}

func DepthDebug(depth int, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(_DEFAULT_DEPTH+depth, DEBUG_LEVEL, args...)
	}
}

func DepthDebugf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(_DEFAULT_DEPTH+depth, DEBUG_LEVEL, format, args...)
	}
}

func Info(args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(_DEFAULT_DEPTH, INFO_LEVEL, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(_DEFAULT_DEPTH, INFO_LEVEL, format, args...)
	}
}

func DepthInfo(depth int, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(_DEFAULT_DEPTH+depth, INFO_LEVEL, args...)
	}
}

func DepthInfof(depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(_DEFAULT_DEPTH+depth, INFO_LEVEL, format, args...)
	}
}

func Warn(args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(_DEFAULT_DEPTH, WARN_LEVEL, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(_DEFAULT_DEPTH, WARN_LEVEL, format, args...)
	}
}

func DepthWarn(depth int, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(_DEFAULT_DEPTH+depth, WARN_LEVEL, args...)
	}
}

func DepthWarnf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(_DEFAULT_DEPTH+depth, WARN_LEVEL, format, args...)
	}
}

func Error(args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(_DEFAULT_DEPTH, ERROR_LEVEL, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(_DEFAULT_DEPTH, ERROR_LEVEL, format, args...)
	}
}

func DepthError(depth int, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(_DEFAULT_DEPTH+depth, ERROR_LEVEL, args...)
	}
}

func DepthErrorf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(_DEFAULT_DEPTH+depth, ERROR_LEVEL, format, args...)
	}
}

func Recover(recover interface{}) {
	defaultLogger.Recover(recover)
}

type logger struct {
	io.Writer
	Header
	enableDebug bool
	enableInfo  bool
	enableWarn  bool
	enableError bool
}

func (lg *logger) SetWriter(writer io.Writer) {
	lg.Writer = writer
}

func (lg *logger) SetHeader(header Header) {
	lg.Header = header
}

func (lg *logger) SetLevel(levels ...Level) {
	lg.enableDebug = false
	lg.enableInfo = false
	lg.enableWarn = false
	lg.enableError = false
	for _, t := range levels {
		switch t {
		case DEBUG_LEVEL:
			lg.enableDebug = true
		case INFO_LEVEL:
			lg.enableInfo = true
		case WARN_LEVEL:
			lg.enableWarn = true
		case ERROR_LEVEL:
			lg.enableDebug = true
		}
	}
}

func (lg *logger) output(depth int, level Level, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, level, depth)
	_log.line = append(_log.line, ' ')
	// log
	fmt.Fprint(_log, args...)
	// wrap
	_log.line = append(_log.line, '\n')
	// output
	lg.Writer.Write(_log.line)
	//
	logPool.Put(_log)
}

func (lg *logger) outputf(depth int, level Level, format string, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, level, depth)
	_log.line = append(_log.line, ' ')
	// log
	fmt.Fprintf(_log, format, args...)
	// wrap
	_log.line = append(_log.line, '\n')
	// output
	lg.Writer.Write(_log.line)
	//
	logPool.Put(_log)
}

func (lg *logger) Debug(args ...interface{}) {
	if lg.enableDebug {
		lg.output(_LOGGER_DEPTH, DEBUG_LEVEL, args...)
	}
}

func (lg *logger) Debugf(format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(_LOGGER_DEPTH, DEBUG_LEVEL, format, args...)
	}
}

func (lg *logger) DepthDebug(depth int, args ...interface{}) {
	if lg.enableDebug {
		lg.output(_LOGGER_DEPTH+depth, DEBUG_LEVEL, args...)
	}
}

func (lg *logger) DepthDebugf(depth int, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(_LOGGER_DEPTH+depth, DEBUG_LEVEL, format, args...)
	}
}

func (lg *logger) Info(args ...interface{}) {
	if lg.enableInfo {
		lg.output(_LOGGER_DEPTH, INFO_LEVEL, args...)
	}
}

func (lg *logger) Infof(format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(_LOGGER_DEPTH, INFO_LEVEL, format, args...)
	}
}

func (lg *logger) DepthInfo(depth int, args ...interface{}) {
	if lg.enableInfo {
		lg.output(_LOGGER_DEPTH+depth, INFO_LEVEL, args...)
	}
}

func (lg *logger) DepthInfof(depth int, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(_LOGGER_DEPTH+depth, INFO_LEVEL, format, args...)
	}
}

func (lg *logger) Warn(args ...interface{}) {
	if lg.enableWarn {
		lg.output(_LOGGER_DEPTH, WARN_LEVEL, args...)
	}
}

func (lg *logger) Warnf(format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(_LOGGER_DEPTH, WARN_LEVEL, format, args...)
	}
}

func (lg *logger) DepthWarn(depth int, args ...interface{}) {
	if lg.enableWarn {
		lg.output(_LOGGER_DEPTH+depth, WARN_LEVEL, args...)
	}
}

func (lg *logger) DepthWarnf(depth int, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(_LOGGER_DEPTH+depth, WARN_LEVEL, format, args...)
	}
}

func (lg *logger) Error(args ...interface{}) {
	if lg.enableError {
		lg.output(_LOGGER_DEPTH, ERROR_LEVEL, args...)
	}
}

func (lg *logger) Errorf(format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(_LOGGER_DEPTH, ERROR_LEVEL, format, args...)
	}
}

func (lg *logger) DepthError(depth int, args ...interface{}) {
	if lg.enableError {
		lg.output(_LOGGER_DEPTH+depth, ERROR_LEVEL, args...)
	}
}

func (lg *logger) DepthErrorf(depth int, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(_LOGGER_DEPTH+depth, ERROR_LEVEL, format, args...)
	}
}

func (lg *logger) Recover(recover interface{}) {
	b := debug.Stack()
	found := false
	n := 0
	for len(b) > 0 {
		i := bytes.IndexByte(b, '\n')
		if i < 0 {
			return
		}
		if !found {
			found = bytes.Contains(b[:i], runtimePanic)
		} else {
			n++
			// the second line
			if n == 2 {
				// skip '\t'
				b := b[1:i]
				_log := logPool.Get().(*Log)
				_log.line = _log.line[:0]
				_log.line = append(_log.line, byte(PANIC_LEVEL))
				_log.line = append(_log.line, ' ')
				// stack
				i = bytes.LastIndexByte(b, ' ')
				if i > 0 {
					b = b[:i]
					// split path and line
					i = bytes.LastIndexByte(b, ':')
					if i > 0 {
						lg.Header.FormatWith(_log, PANIC_LEVEL, string(b[:i]), string(b[i+1:]))
					} else {
						lg.Header.FormatWith(_log, PANIC_LEVEL, "?", "-1")
					}
				} else {
					lg.Header.FormatWith(_log, PANIC_LEVEL, "?", "-1")
				}
				//
				_log.line = append(_log.line, ' ')
				fmt.Fprint(_log, recover)
				_log.line = append(_log.line, '\n')
				lg.Writer.Write(_log.line)
				logPool.Put(_log)
				return
			}
		}
		b = b[i+1:]
	}
}
