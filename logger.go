package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

type OutputType string

const (
	OutputTypeDebug OutputType = "Debug"
	OutputTypeInfo  OutputType = "Info"
	OutputTypeWarn  OutputType = "Warn"
	OutputTypeError OutputType = "Error"
)

const (
	_Debug byte = 'D'
	_Info  byte = 'I'
	_Warn  byte = 'W'
	_Error byte = 'E'
)

const (
	_DefaultDepth = 3
	_LoggerDepth  = 3
)

type Logger interface {
	SetHeader(header Header)
	SetWriter(writer io.Writer)
	SetLevel(level int)
	SetType(types ...OutputType)
	//
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	//
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	//
	DepthDebug(depth int, args ...interface{})
	DepthInfo(depth int, args ...interface{})
	DepthWarn(depth int, args ...interface{})
	DepthError(depth int, args ...interface{})
	//
	DepthDebugf(depth int, format string, args ...interface{})
	DepthInfof(depth int, format string, args ...interface{})
	DepthWarnf(depth int, format string, args ...interface{})
	DepthErrorf(depth int, format string, args ...interface{})
	//
	LevelDebug(level int, args ...interface{})
	LevelInfo(level int, args ...interface{})
	LevelWarn(level int, args ...interface{})
	LevelError(level int, args ...interface{})
	//
	LevelDebugf(level int, format string, args ...interface{})
	LevelInfof(level int, format string, args ...interface{})
	LevelWarnf(level int, format string, args ...interface{})
	LevelErrorf(level int, format string, args ...interface{})
	//
	LevelDepthDebug(level, depth int, args ...interface{})
	LevelDepthInfo(level, depth int, args ...interface{})
	LevelDepthWarn(level, depth int, args ...interface{})
	LevelDepthError(level, depth int, args ...interface{})
	//
	LevelDepthDebugf(level, depth int, format string, args ...interface{})
	LevelDepthInfof(level, depth int, format string, args ...interface{})
	LevelDepthWarnf(level, depth int, format string, args ...interface{})
	LevelDepthErrorf(level, depth int, format string, args ...interface{})
}

var (
	defaultLogger = &logger{
		Writer:      os.Stdout,
		Header:      &CallStackFilePathHeader{},
		level:       0,
		enableDebug: true,
		enableInfo:  true,
		enableWarn:  true,
		enableError: true,
	}
)

func NewLogger(writer io.Writer, level int, header Header, types ...OutputType) Logger {
	lg := new(logger)
	lg.Writer = writer
	lg.Header = header
	lg.level = int32(level)
	lg.SetType(OutputTypeDebug, OutputTypeInfo, OutputTypeWarn, OutputTypeError)
	return lg
}

func SetWriter(writer io.Writer) {
	defaultLogger.SetWriter(writer)
}

func SetHeader(header Header) {
	defaultLogger.SetHeader(header)
}

func SetLevel(level int) {
	defaultLogger.SetLevel(level)
}

func SetType(types ...OutputType) {
	defaultLogger.SetType(types...)
}

func Debug(args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(_DefaultDepth, _Debug, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(_DefaultDepth, _Debug, format, args...)
	}
}

func DepthDebug(depth int, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.output(_DefaultDepth+depth, _Debug, args...)
	}
}

func DepthDebugf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug {
		defaultLogger.outputf(_DefaultDepth+depth, _Debug, format, args...)
	}
}

func Info(args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(_DefaultDepth, _Info, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(_DefaultDepth, _Info, format, args...)
	}
}

func DepthInfo(depth int, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.output(_DefaultDepth+depth, _Info, args...)
	}
}

func DepthInfof(depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo {
		defaultLogger.outputf(_DefaultDepth+depth, _Info, format, args...)
	}
}

func Warn(args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(_DefaultDepth, _Warn, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(_DefaultDepth, _Warn, format, args...)
	}
}

func DepthWarn(depth int, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.output(_DefaultDepth+depth, _Warn, args...)
	}
}

func DepthWarnf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn {
		defaultLogger.outputf(_DefaultDepth+depth, _Warn, format, args...)
	}
}

func Error(args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(_DefaultDepth, _Error, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(_DefaultDepth, _Error, format, args...)
	}
}

func DepthError(depth int, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.output(_DefaultDepth+depth, _Error, args...)
	}
}

func DepthErrorf(depth int, format string, args ...interface{}) {
	if defaultLogger.enableError {
		defaultLogger.outputf(_DefaultDepth+depth, _Error, format, args...)
	}
}

func LevelDebug(level int, args ...interface{}) {
	if defaultLogger.enableDebug && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth, _Debug, args...)
	}
}

func LevelDebugf(level int, format string, args ...interface{}) {
	if defaultLogger.enableDebug && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth, _Debug, format, args...)
	}
}

func LevelDepthDebug(level int, depth int, args ...interface{}) {
	if defaultLogger.enableDebug && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth+depth, _Debug, args...)
	}
}

func LevelDepthDebugf(level int, depth int, format string, args ...interface{}) {
	if defaultLogger.enableDebug && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth+depth, _Debug, format, args...)
	}
}

func LevelInfo(level int, args ...interface{}) {
	if defaultLogger.enableInfo && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth, _Info, args...)
	}
}

func LevelInfof(level int, format string, args ...interface{}) {
	if defaultLogger.enableInfo && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth, _Info, format, args...)
	}
}

func LevelDepthInfo(level int, depth int, args ...interface{}) {
	if defaultLogger.enableInfo && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth+depth, _Info, args...)
	}
}

func LevelDepthInfof(level int, depth int, format string, args ...interface{}) {
	if defaultLogger.enableInfo && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth+depth, _Info, format, args...)
	}
}

func LevelWarn(level int, args ...interface{}) {
	if defaultLogger.enableWarn && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth, _Warn, args...)
	}
}

func LevelWarnf(level int, format string, args ...interface{}) {
	if defaultLogger.enableWarn && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth, _Warn, format, args...)
	}
}

func LevelDepthWarn(level int, depth int, args ...interface{}) {
	if defaultLogger.enableWarn && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth+depth, _Warn, args...)
	}
}

func LevelDepthWarnf(level int, depth int, format string, args ...interface{}) {
	if defaultLogger.enableWarn && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth+depth, _Warn, format, args...)
	}
}

func LevelError(level int, args ...interface{}) {
	if defaultLogger.enableError && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth, _Error, args...)
	}
}

func LevelErrorf(level int, format string, args ...interface{}) {
	if defaultLogger.enableError && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth, _Error, format, args...)
	}
}

func LevelDepthError(level int, depth int, args ...interface{}) {
	if defaultLogger.enableError && defaultLogger.level <= int32(level) {
		defaultLogger.output(_DefaultDepth+depth, _Error, args...)
	}
}

func LevelDepthErrorf(level int, depth int, format string, args ...interface{}) {
	if defaultLogger.enableError && defaultLogger.level <= int32(level) {
		defaultLogger.outputf(_DefaultDepth+depth, _Error, format, args...)
	}
}

type logger struct {
	io.Writer
	Header
	level       int32
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

func (lg *logger) SetLevel(level int) {
	atomic.SwapInt32(&lg.level, int32(level))
}

func (lg *logger) SetType(types ...OutputType) {
	lg.enableDebug = false
	lg.enableInfo = false
	lg.enableWarn = false
	lg.enableError = false
	for _, t := range types {
		switch t {
		case OutputTypeDebug:
			lg.enableDebug = true
		case OutputTypeInfo:
			lg.enableInfo = true
		case OutputTypeWarn:
			lg.enableWarn = true
		case OutputTypeError:
			lg.enableDebug = true
		}
	}
}

func (lg *logger) output(depth int, _type byte, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	_log.line = append(_log.line, _type)
	_log.line = append(_log.line, ' ')
	lg.Header.Format(_log, depth)
	_log.line = append(_log.line, ' ')
	fmt.Fprint(_log, args...)
	_log.line = append(_log.line, '\n')
	lg.Writer.Write(_log.line)
	logPool.Put(_log)
}

func (lg *logger) outputf(depth int, _type byte, format string, args ...interface{}) {
	_log := logPool.Get().(*Log)
	_log.line = _log.line[:0]
	_log.line = append(_log.line, _type)
	_log.line = append(_log.line, ' ')
	lg.Header.Format(_log, depth)
	_log.line = append(_log.line, ' ')
	fmt.Fprintf(_log, format, args...)
	_log.line = append(_log.line, '\n')
	lg.Writer.Write(_log.line)
	logPool.Put(_log)
}

func (lg *logger) Debug(args ...interface{}) {
	if lg.enableDebug {
		lg.output(_LoggerDepth, _Debug, args...)
	}
}

func (lg *logger) Debugf(format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(_LoggerDepth, _Debug, format, args...)
	}
}

func (lg *logger) DepthDebug(depth int, args ...interface{}) {
	if lg.enableDebug {
		lg.output(_LoggerDepth+depth, _Debug, args...)
	}
}

func (lg *logger) DepthDebugf(depth int, format string, args ...interface{}) {
	if lg.enableDebug {
		lg.outputf(_LoggerDepth+depth, _Debug, format, args...)
	}
}

func (lg *logger) Info(args ...interface{}) {
	if lg.enableInfo {
		lg.output(_LoggerDepth, _Info, args...)
	}
}

func (lg *logger) Infof(format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(_LoggerDepth, _Info, format, args...)
	}
}

func (lg *logger) DepthInfo(depth int, args ...interface{}) {
	if lg.enableInfo {
		lg.output(_LoggerDepth+depth, _Info, args...)
	}
}

func (lg *logger) DepthInfof(depth int, format string, args ...interface{}) {
	if lg.enableInfo {
		lg.outputf(_LoggerDepth+depth, _Info, format, args...)
	}
}

func (lg *logger) Warn(args ...interface{}) {
	if lg.enableWarn {
		lg.output(_LoggerDepth, _Warn, args...)
	}
}

func (lg *logger) Warnf(format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(_LoggerDepth, _Warn, format, args...)
	}
}

func (lg *logger) DepthWarn(depth int, args ...interface{}) {
	if lg.enableWarn {
		lg.output(_LoggerDepth+depth, _Warn, args...)
	}
}

func (lg *logger) DepthWarnf(depth int, format string, args ...interface{}) {
	if lg.enableWarn {
		lg.outputf(_LoggerDepth+depth, _Warn, format, args...)
	}
}

func (lg *logger) Error(args ...interface{}) {
	if lg.enableError {
		lg.output(_LoggerDepth, _Error, args...)
	}
}

func (lg *logger) Errorf(format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(_LoggerDepth, _Error, format, args...)
	}
}

func (lg *logger) DepthError(depth int, args ...interface{}) {
	if lg.enableError {
		lg.output(_LoggerDepth+depth, _Error, args...)
	}
}

func (lg *logger) DepthErrorf(depth int, format string, args ...interface{}) {
	if lg.enableError {
		lg.outputf(_LoggerDepth+depth, _Error, format, args...)
	}
}

func (lg *logger) LevelDebug(level int, args ...interface{}) {
	if lg.enableDebug && lg.level <= int32(level) {
		lg.output(_LoggerDepth, _Debug, args...)
	}
}

func (lg *logger) LevelDebugf(level int, format string, args ...interface{}) {
	if lg.enableDebug && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth, _Debug, format, args...)
	}
}

func (lg *logger) LevelDepthDebug(level, depth int, args ...interface{}) {
	if lg.enableDebug && lg.level <= int32(level) {
		lg.output(_LoggerDepth+depth, _Debug, args...)
	}
}

func (lg *logger) LevelDepthDebugf(level, depth int, format string, args ...interface{}) {
	if lg.enableDebug && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth+depth, _Debug, format, args...)
	}
}

func (lg *logger) LevelInfo(level int, args ...interface{}) {
	if lg.enableInfo && lg.level <= int32(level) {
		lg.output(_LoggerDepth, _Info, args...)
	}
}

func (lg *logger) LevelInfof(level int, format string, args ...interface{}) {
	if lg.enableInfo && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth, _Info, format, args...)
	}
}

func (lg *logger) LevelDepthInfo(level, depth int, args ...interface{}) {
	if lg.enableInfo && lg.level <= int32(level) {
		lg.output(depth, _Info, args...)
	}
}

func (lg *logger) LevelDepthInfof(level, depth int, format string, args ...interface{}) {
	if lg.enableInfo && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth+depth, _Info, format, args...)
	}
}

func (lg *logger) LevelWarn(level int, args ...interface{}) {
	if lg.enableWarn && lg.level <= int32(level) {
		lg.output(_LoggerDepth, _Warn, args...)
	}
}

func (lg *logger) LevelWarnf(level int, format string, args ...interface{}) {
	if lg.enableWarn && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth, _Warn, format, args...)
	}
}

func (lg *logger) LevelDepthWarn(level, depth int, args ...interface{}) {
	if lg.enableWarn && lg.level <= int32(level) {
		lg.output(_LoggerDepth+depth, _Warn, args...)
	}
}

func (lg *logger) LevelDepthWarnf(level, depth int, format string, args ...interface{}) {
	if lg.enableWarn && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth+depth, _Warn, format, args...)
	}
}

func (lg *logger) LevelError(level int, args ...interface{}) {
	if lg.enableError && lg.level <= int32(level) {
		lg.output(_LoggerDepth, _Error, args...)
	}
}

func (lg *logger) LevelErrorf(level int, format string, args ...interface{}) {
	if lg.enableError && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth, _Error, format, args...)
	}
}

func (lg *logger) LevelDepthError(level, depth int, args ...interface{}) {
	if lg.enableError && lg.level <= int32(level) {
		lg.output(_LoggerDepth+depth, _Error, args...)
	}
}

func (lg *logger) LevelDepthErrorf(level, depth int, format string, args ...interface{}) {
	if lg.enableError && lg.level <= int32(level) {
		lg.outputf(_LoggerDepth+depth, _Error, format, args...)
	}
}
