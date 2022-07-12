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
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugDepth(depth int, args ...interface{})
	DebugfDepth(depth int, format string, args ...interface{})
	DebugTrace(traceID string, args ...interface{})
	DebugfTrace(traceID string, format string, args ...interface{})
	DebugDepthTrace(traceID string, depth int, args ...interface{})
	DebugfDepthTrace(traceID string, depth int, format string, args ...interface{})
	// Info 级别的方法。
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoDepth(depth int, args ...interface{})
	InfofDepth(depth int, format string, args ...interface{})
	InfoTrace(traceID string, args ...interface{})
	InfofTrace(traceID string, format string, args ...interface{})
	InfoDepthTrace(traceID string, depth int, args ...interface{})
	InfofDepthTrace(traceID string, depth int, format string, args ...interface{})
	// Warn 级别的方法。
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnDepth(depth int, args ...interface{})
	WarnfDepth(depth int, format string, args ...interface{})
	WarnTrace(traceID string, args ...interface{})
	WarnfTrace(traceID string, format string, args ...interface{})
	WarnDepthTrace(traceID string, depth int, args ...interface{})
	WarnfDepthTrace(traceID string, depth int, format string, args ...interface{})
	// Error 级别的方法。
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorDepth(depth int, args ...interface{})
	ErrorfDepth(depth int, format string, args ...interface{})
	ErrorTrace(traceID string, args ...interface{})
	ErrorfTrace(traceID string, format string, args ...interface{})
	ErrorDepthTrace(traceID string, depth int, args ...interface{})
	ErrorfDepthTrace(traceID string, depth int, format string, args ...interface{})
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
	// 使用网卡来初始化默认 logger 的 appID 。
	addr, err := net.Interfaces()
	if nil != err {
		panic(err)
	}
	// 第一个网卡的 MAC 地址
	appID := ""
	for i := 0; i < len(addr); i++ {
		if addr[i].Flags|net.FlagUp != 0 && len(addr[i].HardwareAddr) != 0 {
			appID = addr[i].HardwareAddr.String()
			break
		}
	}
	// 默认 logger
	defaultLogger = &logger{
		Writer:      os.Stdout,
		Header:      NewHeaderFormater("filePathStack", appID),
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
func Debug(args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output("", _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// Debugf 使用默认 Logger 输出日志。
func Debugf(format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf("", _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// DebugDepth 使用默认 Logger 输出日志。
func DebugDepth(depth int, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output("", _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// DebugfDepth 使用默认 Logger 输出日志。
func DebugfDepth(depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf("", _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// DebugTrace 使用默认 Logger 输出日志。
func DebugTrace(traceID string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// DebugfTrace 使用默认 Logger 输出日志。
func DebugfTrace(traceID string, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// DebugDepthTrace 使用默认 Logger 输出日志。
func DebugDepthTrace(traceID string, depth int, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, DebugLevel, args...)
	}
}

// DebugfDepthTrace 使用默认 Logger 输出日志。
func DebugfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, DebugLevel, format, args...)
	}
}

// Info 使用默认 Logger 输出日志。
func Info(args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output("", _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// Infof 使用默认 Logger 输出日志。
func Infof(format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf("", _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// InfoDepth 使用默认 Logger 输出日志。
func InfoDepth(depth int, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output("", _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// InfofDepth 使用默认 Logger 输出日志。
func InfofDepth(depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf("", _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// InfoTrace 使用默认 Logger 输出日志。
func InfoTrace(traceID string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// InfofTrace 使用默认 Logger 输出日志。
func InfofTrace(traceID string, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// InfoDepthTrace 使用默认 Logger 输出日志。
func InfoDepthTrace(traceID string, depth int, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, InfoLevel, args...)
	}
}

// InfofDepthTrace 使用默认 Logger 输出日志。
func InfofDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, InfoLevel, format, args...)
	}
}

// Warn 使用默认 Logger 输出日志。
func Warn(args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output("", _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// Warnf 使用默认 Logger 输出日志。
func Warnf(format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf("", _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// WarnDepth 使用默认 Logger 输出日志。
func WarnDepth(depth int, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output("", _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// WarnfDepth 使用默认 Logger 输出日志。
func WarnfDepth(depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf("", _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// WarnTrace 使用默认 Logger 输出日志。
func WarnTrace(traceID string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// WarnfTrace 使用默认 Logger 输出日志。
func WarnfTrace(traceID string, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// WarnDepthTrace 使用默认 Logger 输出日志。
func WarnDepthTrace(traceID string, depth int, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, WarnLevel, args...)
	}
}

// WarnfDepthTrace 使用默认 Logger 输出日志。
func WarnfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, WarnLevel, format, args...)
	}
}

// Error 使用默认 Logger 输出日志。
func Error(args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output("", _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// Errorf 使用默认 Logger 输出日志。
func Errorf(format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf("", _DEFAULT_DEPTH, ErrorLevel, format, args...)
	}
}

// ErrorDepth 使用默认 Logger 输出日志。
func ErrorDepth(depth int, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output("", _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// ErrorfDepth 使用默认 Logger 输出日志。
func ErrorfDepth(depth int, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf("", _DEFAULT_DEPTH, ErrorLevel, format, args...)
	}
}

// Error 使用默认 Logger 输出日志。
func ErrorTrace(traceID string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// Errorf 使用默认 Logger 输出日志。
func ErrorfTrace(traceID string, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, ErrorLevel, format, args...)
	}
}

// ErrorDepthTrace 使用默认 Logger 输出日志。
func ErrorDepthTrace(traceID string, depth int, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(traceID, _DEFAULT_DEPTH, ErrorLevel, args...)
	}
}

// ErrorfDepthTrace 使用默认 Logger 输出日志。
func ErrorfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(traceID, _DEFAULT_DEPTH, ErrorLevel, format, args...)
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

func (lg *logger) output(traceID string, depth int, level Level, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, traceID, level, depth)
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

func (lg *logger) outputf(traceID string, depth int, level Level, format string, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	// header
	lg.Header.Format(_log, traceID, level, depth)
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

func (lg *logger) Debug(args ...interface{}) {
	if lg.enableDebug {
		lg.output("", _LOGGER_DEPTH, DebugLevel, args...)
	}
}

func (lg *logger) Debugf(format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf("", _LOGGER_DEPTH, DebugLevel, format, args...)
	}
}

func (lg *logger) DebugDepth(depth int, args ...interface{}) {
	if lg.enableDebug {
		lg.output("", _LOGGER_DEPTH+depth, DebugLevel, args...)
	}
}

func (lg *logger) DebugfDepth(depth int, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf("", _LOGGER_DEPTH+depth, DebugLevel, format, args...)
	}
}

func (lg *logger) DebugTrace(traceID string, args ...interface{}) {
	if lg.enableDebug {
		lg.output(traceID, _LOGGER_DEPTH, DebugLevel, args...)
	}
}

func (lg *logger) DebugfTrace(traceID string, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(traceID, _LOGGER_DEPTH, DebugLevel, format, args...)
	}
}

func (lg *logger) DebugDepthTrace(traceID string, depth int, args ...interface{}) {
	if lg.enableDebug {
		lg.output(traceID, _LOGGER_DEPTH+depth, DebugLevel, args...)
	}
}

func (lg *logger) DebugfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(traceID, _LOGGER_DEPTH+depth, DebugLevel, format, args...)
	}
}

func (lg *logger) Info(args ...interface{}) {
	if lg.enableInfo {
		lg.output("", _LOGGER_DEPTH, InfoLevel, args...)
	}
}

func (lg *logger) Infof(format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf("", _LOGGER_DEPTH, InfoLevel, format, args...)
	}
}

func (lg *logger) InfoDepth(depth int, args ...interface{}) {
	if lg.enableInfo {
		lg.output("", _LOGGER_DEPTH+depth, InfoLevel, args...)
	}
}

func (lg *logger) InfofDepth(depth int, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf("", _LOGGER_DEPTH+depth, InfoLevel, format, args...)
	}
}

func (lg *logger) InfoTrace(traceID string, args ...interface{}) {
	if lg.enableInfo {
		lg.output(traceID, _LOGGER_DEPTH, InfoLevel, args...)
	}
}

func (lg *logger) InfofTrace(traceID string, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(traceID, _LOGGER_DEPTH, InfoLevel, format, args...)
	}
}

func (lg *logger) InfoDepthTrace(traceID string, depth int, args ...interface{}) {
	if lg.enableInfo {
		lg.output(traceID, _LOGGER_DEPTH+depth, InfoLevel, args...)
	}
}

func (lg *logger) InfofDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(traceID, _LOGGER_DEPTH+depth, InfoLevel, format, args...)
	}
}

func (lg *logger) Warn(args ...interface{}) {
	if lg.enableWarn {
		lg.output("", _LOGGER_DEPTH, WarnLevel, args...)
	}
}

func (lg *logger) Warnf(format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf("", _LOGGER_DEPTH, WarnLevel, format, args...)
	}
}

func (lg *logger) WarnDepth(depth int, args ...interface{}) {
	if lg.enableWarn {
		lg.output("", _LOGGER_DEPTH+depth, WarnLevel, args...)
	}
}

func (lg *logger) WarnfDepth(depth int, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf("", _LOGGER_DEPTH+depth, WarnLevel, format, args...)
	}
}

func (lg *logger) WarnTrace(traceID string, args ...interface{}) {
	if lg.enableWarn {
		lg.output(traceID, _LOGGER_DEPTH, WarnLevel, args...)
	}
}

func (lg *logger) WarnfTrace(traceID string, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(traceID, _LOGGER_DEPTH, WarnLevel, format, args...)
	}
}

func (lg *logger) WarnDepthTrace(traceID string, depth int, args ...interface{}) {
	if lg.enableWarn {
		lg.output(traceID, _LOGGER_DEPTH+depth, WarnLevel, args...)
	}
}

func (lg *logger) WarnfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(traceID, _LOGGER_DEPTH+depth, WarnLevel, format, args...)
	}
}

func (lg *logger) Error(args ...interface{}) {
	if lg.enableError {
		lg.output("", _LOGGER_DEPTH, ErrorLevel, args...)
	}
}

func (lg *logger) Errorf(format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf("", _LOGGER_DEPTH, ErrorLevel, format, args...)
	}
}

func (lg *logger) ErrorDepth(depth int, args ...interface{}) {
	if lg.enableError {
		lg.output("", _LOGGER_DEPTH+depth, ErrorLevel, args...)
	}
}

func (lg *logger) ErrorfDepth(depth int, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf("", _LOGGER_DEPTH+depth, ErrorLevel, format, args...)
	}
}

func (lg *logger) ErrorTrace(traceID string, args ...interface{}) {
	if lg.enableError {
		lg.output(traceID, _LOGGER_DEPTH, ErrorLevel, args...)
	}
}

func (lg *logger) ErrorfTrace(traceID string, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(traceID, _LOGGER_DEPTH, ErrorLevel, format, args...)
	}
}

func (lg *logger) ErrorDepthTrace(traceID string, depth int, args ...interface{}) {
	if lg.enableError {
		lg.output(traceID, _LOGGER_DEPTH+depth, ErrorLevel, args...)
	}
}

func (lg *logger) ErrorfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(traceID, _LOGGER_DEPTH+depth, ErrorLevel, format, args...)
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
