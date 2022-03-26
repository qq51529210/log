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
	DebugLevel Level = 'D'
	InfoLevel  Level = 'I'
	WarnLevel  Level = 'W'
	ErrorLevel Level = 'E'
	// 这个级别在 recover panic 的时候自动使用的。
	PanicLevel Level = 'P'
)

const (
	// 默认日志的 depth
	_DEFAULT_DEPTH = 3
	// 日志的默认 depth
	_LOGGER_DEPTH = 3
)

// Logger 表示一个日志记录器。
type Logger interface {
	// 设置日志头格式化接口。
	SetHeaderFormater(headerFormater HeaderFormater)
	// 设置日志的输出。
	SetOutput(output io.Writer)
	// 设置允许输出的日志的级别，没有被设置的级别不会输出日志，一般在程序运行的时候设置。
	SetLevel(levels ...Level)
	// Debug 级别的方法。
	Debug(trackId string, args ...interface{})
	Debugf(trackId string, format string, args ...interface{})
	DepthDebug(trackId string, depth int, args ...interface{})
	DepthDebugf(trackId string, depth int, format string, args ...interface{})
	// Info 级别的方法。
	Info(trackId string, args ...interface{})
	Infof(trackId string, format string, args ...interface{})
	DepthInfo(trackId string, depth int, args ...interface{})
	DepthInfof(trackId string, depth int, format string, args ...interface{})
	// Warn 级别的方法。
	Warn(trackId string, args ...interface{})
	Warnf(trackId string, format string, args ...interface{})
	DepthWarn(trackId string, depth int, args ...interface{})
	DepthWarnf(trackId string, depth int, format string, args ...interface{})
	// Error 级别的方法。
	Error(trackId string, args ...interface{})
	Errorf(trackId string, format string, args ...interface{})
	DepthError(trackId string, depth int, args ...interface{})
	DepthErrorf(trackId string, depth int, format string, args ...interface{})
	// Recover 检索出 panic 的文件和行数，如果 recover 不为 nil 。
	Recover(recover interface{})
}

var (
	// 默认 logger ，在 init 函数中初始化。
	defaultLogger *logger
	// 用于检索 panic 堆栈那一行的信息。
	runtimePanic = []byte("runtime/panic.go")
)

func init() {
	// 使用网卡来初始化默认 logger 的 appId 。
	addr, err := net.Interfaces()
	if nil != err {
		panic(err)
	}
	// 第一个网卡的 MAC 地址
	appId := ""
	for i := 0; i < len(addr); i++ {
		if addr[i].Flags|net.FlagUp != 0 && len(addr[i].HardwareAddr) != 0 {
			appId = addr[i].HardwareAddr.String()
			break
		}
	}
	// 默认 logger
	defaultLogger = &logger{
		Writer:      os.Stdout,
		Header:      NewHeaderFormater("filePathStack", appId),
		enableDebug: true,
		enableInfo:  true,
		enableWarn:  true,
		enableError: true,
	}
}

// NewLogger 返回一个 Logger 实例。output 和 headerFormater 是初始化参数。
func NewLogger(output io.Writer, headerFormater HeaderFormater) Logger {
	lg := new(logger)
	lg.Writer = output
	lg.Header = headerFormater
	lg.enableDebug = true
	lg.enableInfo = true
	lg.enableWarn = true
	lg.enableError = true
	return lg
}

// SetOutput 设置默认 Logger 的输出。
func SetOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

// SetHeaderFormater 设置默认 Logger 的 HeaderFormater 。
func SetHeaderFormater(headerFormater HeaderFormater) {
	defaultLogger.SetHeaderFormater(headerFormater)
}

// SetLevel 设置默认 Logger 的输出级别 。
func SetLevel(levels ...Level) {
	defaultLogger.SetLevel(levels...)
}

// Debug 使用默认 Logger 输出日志。
func Debug(trackId string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// Debugf 使用默认 Logger 输出日志。
func Debugf(trackId string, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// DepthDebug 使用默认 Logger 输出日志。
func DepthDebug(trackId string, depth int, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// DepthDebugf 使用默认 Logger 输出日志。
func DepthDebugf(trackId string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// Info 使用默认 Logger 输出日志。
func Info(trackId string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// Infof 使用默认 Logger 输出日志。
func Infof(trackId string, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// DepthInfo 使用默认 Logger 输出日志。
func DepthInfo(trackId string, depth int, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// DepthInfof 使用默认 Logger 输出日志。
func DepthInfof(trackId string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// Warn 使用默认 Logger 输出日志。
func Warn(trackId string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// Warnf 使用默认 Logger 输出日志。
func Warnf(trackId string, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// DepthWarn 使用默认 Logger 输出日志。
func DepthWarn(trackId string, depth int, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// DepthWarnf 使用默认 Logger 输出日志。
func DepthWarnf(trackId string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// Error 使用默认 Logger 输出日志。
func Error(trackId string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// Errorf 使用默认 Logger 输出日志。
func Errorf(trackId string, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, ErrorLevel, format, args...)
	}
}

// DepthError 使用默认 Logger 输出日志。
func DepthError(trackId string, depth int, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(trackId, _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// DepthErrorf 使用默认 Logger 输出日志。
func DepthErrorf(trackId string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(trackId, _DEFAULT_DEPTH, ErrorLevel, format, args...)
	}
}

// Recover 使用默认 Logger 输出日志。
func Recover(recover interface{}) {
	defaultLogger.Recover(recover)
}

type logger struct {
	io.Writer
	Header      HeaderFormater
	enableDebug bool
	enableInfo  bool
	enableWarn  bool
	enableError bool
}

func (lg *logger) output(trackId string, depth int, level Level, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, trackId, level, depth)
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

func (lg *logger) outputf(trackId string, depth int, level Level, format string, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, trackId, level, depth)
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

func (lg *logger) SetOutput(output io.Writer) {
	lg.Writer = output
}

func (lg *logger) SetHeaderFormater(headerFormater HeaderFormater) {
	lg.Header = headerFormater
}

func (lg *logger) SetLevel(levels ...Level) {
	lg.enableDebug = false
	lg.enableInfo = false
	lg.enableWarn = false
	lg.enableError = false
	for _, t := range levels {
		switch t {
		case DebugLevel:
			lg.enableDebug = true
		case InfoLevel:
			lg.enableInfo = true
		case WarnLevel:
			lg.enableWarn = true
		case ErrorLevel:
			lg.enableDebug = true
		}
	}
}

func (lg *logger) Debug(trackId string, args ...interface{}) {
	if lg.enableDebug {
		lg.output(trackId, _LOGGER_DEPTH, DebugLevel, args...)
	}
}

func (lg *logger) Debugf(trackId string, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(trackId, _LOGGER_DEPTH, DebugLevel, format, args...)
	}
}

func (lg *logger) DepthDebug(trackId string, depth int, args ...interface{}) {
	if lg.enableDebug {
		lg.output(trackId, _LOGGER_DEPTH+depth, DebugLevel, args...)
	}
}

func (lg *logger) DepthDebugf(trackId string, depth int, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(trackId, _LOGGER_DEPTH+depth, DebugLevel, format, args...)
	}
}

func (lg *logger) Info(trackId string, args ...interface{}) {
	if lg.enableInfo {
		lg.output(trackId, _LOGGER_DEPTH, InfoLevel, args...)
	}
}

func (lg *logger) Infof(trackId string, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(trackId, _LOGGER_DEPTH, InfoLevel, format, args...)
	}
}

func (lg *logger) DepthInfo(trackId string, depth int, args ...interface{}) {
	if lg.enableInfo {
		lg.output(trackId, _LOGGER_DEPTH+depth, InfoLevel, args...)
	}
}

func (lg *logger) DepthInfof(trackId string, depth int, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(trackId, _LOGGER_DEPTH+depth, InfoLevel, format, args...)
	}
}

func (lg *logger) Warn(trackId string, args ...interface{}) {
	if lg.enableWarn {
		lg.output(trackId, _LOGGER_DEPTH, WarnLevel, args...)
	}
}

func (lg *logger) Warnf(trackId string, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(trackId, _LOGGER_DEPTH, WarnLevel, format, args...)
	}
}

func (lg *logger) DepthWarn(trackId string, depth int, args ...interface{}) {
	if lg.enableWarn {
		lg.output(trackId, _LOGGER_DEPTH+depth, WarnLevel, args...)
	}
}

func (lg *logger) DepthWarnf(trackId string, depth int, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(trackId, _LOGGER_DEPTH+depth, WarnLevel, format, args...)
	}
}

func (lg *logger) Error(trackId string, args ...interface{}) {
	if lg.enableError {
		lg.output(trackId, _LOGGER_DEPTH, ErrorLevel, args...)
	}
}

func (lg *logger) Errorf(trackId string, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(trackId, _LOGGER_DEPTH, ErrorLevel, format, args...)
	}
}

func (lg *logger) DepthError(trackId string, depth int, args ...interface{}) {
	if lg.enableError {
		lg.output(trackId, _LOGGER_DEPTH+depth, ErrorLevel, args...)
	}
}

func (lg *logger) DepthErrorf(trackId string, depth int, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(trackId, _LOGGER_DEPTH+depth, ErrorLevel, format, args...)
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
				// stack
				i = bytes.LastIndexByte(b, ' ')
				if i > 0 {
					b = b[:i]
					// split path and line
					i = bytes.LastIndexByte(b, ':')
					if i > 0 {
						lg.Header.FormatWith(_log, "", PanicLevel, string(b[:i]), string(b[i+1:]))
					} else {
						lg.Header.FormatWith(_log, "", PanicLevel, "?", "-1")
					}
				} else {
					lg.Header.FormatWith(_log, "", PanicLevel, "?", "-1")
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
