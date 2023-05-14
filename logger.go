package log

import (
	"fmt"
	"io"
	"runtime"
)

var (
	// header:
	headerEnd = []byte(": ")
)

const (
	loggerDepth = 3
)

// Logger 日志接口
type Logger interface {
	// 如果 recover 不为 nil 则打印堆栈
	Recover(recover any)
	// Debug 级别的方法
	IsDebug() bool
	EnableDebug(enable bool)
	Debug(args ...any)
	Debugf(format string, args ...any)
	DebugDepth(depth int, args ...any)
	DebugfDepth(depth int, format string, args ...any)
	DebugTrace(traceID string, args ...any)
	DebugfTrace(traceID string, format string, args ...any)
	DebugDepthTrace(depth int, traceID string, args ...any)
	DebugfDepthTrace(depth int, traceID string, format string, args ...any)
	// Info 级别的方法。
	IsInfo() bool
	EnableInfo(enable bool)
	Info(args ...any)
	Infof(format string, args ...any)
	InfoDepth(depth int, args ...any)
	InfofDepth(depth int, format string, args ...any)
	InfoTrace(traceID string, args ...any)
	InfofTrace(traceID string, format string, args ...any)
	InfoDepthTrace(depth int, traceID string, args ...any)
	InfofDepthTrace(depth int, traceID string, format string, args ...any)
	// Warn 级别的方法。
	IsWarn() bool
	EnableWarn(enable bool)
	Warn(args ...any)
	Warnf(format string, args ...any)
	WarnDepth(depth int, args ...any)
	WarnfDepth(depth int, format string, args ...any)
	WarnTrace(traceID string, args ...any)
	WarnfTrace(traceID string, format string, args ...any)
	WarnDepthTrace(depth int, traceID string, args ...any)
	WarnfDepthTrace(depth int, traceID string, format string, args ...any)
	// Error 级别的方法。
	IsError() bool
	EnableError(enable bool)
	Error(args ...any)
	Errorf(format string, args ...any)
	ErrorDepth(depth int, args ...any)
	ErrorfDepth(depth int, format string, args ...any)
	ErrorTrace(traceID string, args ...any)
	ErrorfTrace(traceID string, format string, args ...any)
	ErrorDepthTrace(depth int, traceID string, args ...any)
	ErrorfDepthTrace(depth int, traceID string, format string, args ...any)
}

// logger 默认实现
type logger struct {
	// 输出
	io.Writer
	// 头格式
	Header
	// 名称
	name string
	// 是否禁止 debug
	disableDebug bool
	// 是否禁止 info
	disableInfo bool
	// 是否禁止 warn
	disableWarn bool
	// 是否禁止 error
	disableError bool
}

// NewLogger 返回默认的 Logger
// 格式 "Header [name] [level] [tracID] text"
func NewLogger(writer io.Writer, header Header, name string) Logger {
	lg := new(logger)
	lg.Writer = writer
	lg.Header = header
	if name != "" {
		lg.name = fmt.Sprintf("[%s] ", name)
	}
	return lg
}

func (lg *logger) print(depth int, level *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *logger) printf(depth int, level, format *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 日志
	fmt.Fprintf(l, *format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *logger) printTrace(depth int, trace, level *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, *trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *logger) printfTrace(depth int, trace, level, format *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, *trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprintf(l, *format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *logger) Recover(recover any) {
	if recover == nil {
		return
	}
	// get statck info l.line
	b := logPool.Get().(*Log)
	b.b = b.b[:cap(b.b)]
	for {
		n := runtime.Stack(b.b, false)
		if n < len(b.b) {
			b.b = b.b[:n]
			break
		}
		b.b = make([]byte, len(b.b)+1024)
	}
	// log
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, panicLevel...)
	lg.Header.Time(l)
	// recover
	fmt.Fprintf(l, ": %v\n", recover)
	// filter
	p := b.b
	// 找到 panic.go
	found := false
	for len(p) > 0 {
		// find new line
		i := 0
		for ; i < len(p); i++ {
			if p[i] == '\n' {
				i++
				break
			}
		}
		line := p[:i]
		p = p[i:]
		// find file line
		if line[0] != '\t' {
			continue
		}
		if !found {
			found = isPanicGO(line)
			continue
		}
		// \t filepath/file.go:line +0x622
		i = len(line) - 1
		for i > 0 {
			if line[i] == ' ' {
				//
				line = line[:i]
				break
			}
			i--
		}
		// write
		l.b = append(l.b, "[statck] "...)
		l.b = append(l.b, line[1:]...)
		l.b = append(l.b, '\n')
	}
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(b)
	logPool.Put(l)
}

func isPanicGO(line []byte) bool {
	for i := len(line) - 1; i > 1; i-- {
		if line[i] == '/' {
			for j := i; j < len(line); j++ {
				if line[j] == 'p' &&
					line[j+1] == 'a' &&
					line[j+2] == 'n' &&
					line[j+3] == 'i' &&
					line[j+4] == 'c' &&
					line[j+5] == '.' &&
					line[j+6] == 'g' &&
					line[j+7] == 'o' {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (lg *logger) IsDebug() bool {
	return !lg.disableDebug
}

func (lg *logger) EnableDebug(enable bool) {
	lg.disableDebug = !enable
}

func (lg *logger) Debug(args ...any) {
	if lg.disableDebug {
		return
	}
	lg.print(loggerDepth, &debugLevel, args...)
}

func (lg *logger) Debugf(format string, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.printf(loggerDepth, &debugLevel, &format, args...)
}

func (lg *logger) DebugDepth(depth int, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.print(loggerDepth+depth, &debugLevel, args...)
}

func (lg *logger) DebugfDepth(depth int, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.printf(loggerDepth+depth, &debugLevel, &format, args...)
}

func (lg *logger) DebugTrace(traceID string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &debugLevel, args...)
	} else {
		lg.print(loggerDepth, &debugLevel, args...)
	}
}

func (lg *logger) DebugfTrace(traceID, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &debugLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &debugLevel, &format, args...)
	}
}

func (lg *logger) DebugDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &debugLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &debugLevel, args...)
	}
}

func (lg *logger) DebugfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &debugLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &debugLevel, &format, args...)
	}
}

func (lg *logger) IsInfo() bool {
	return !lg.disableInfo
}

func (lg *logger) EnableInfo(enable bool) {
	lg.disableInfo = !enable
}

func (lg *logger) Info(args ...any) {
	if lg.disableInfo {
		return
	}
	lg.print(loggerDepth, &infoLevel, args...)
}

func (lg *logger) Infof(format string, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.printf(loggerDepth, &infoLevel, &format, args...)
}

func (lg *logger) InfoDepth(depth int, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.print(loggerDepth+depth, &infoLevel, args...)
}

func (lg *logger) InfofDepth(depth int, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.printf(loggerDepth+depth, &infoLevel, &format, args...)
}

func (lg *logger) InfoTrace(traceID string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &infoLevel, args...)
	} else {
		lg.print(loggerDepth, &infoLevel, args...)
	}
}

func (lg *logger) InfofTrace(traceID, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &infoLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &infoLevel, &format, args...)
	}
}

func (lg *logger) InfoDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &infoLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &infoLevel, args...)
	}
}

func (lg *logger) InfofDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &infoLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &infoLevel, &format, args...)
	}
}

func (lg *logger) IsWarn() bool {
	return !lg.disableWarn
}

func (lg *logger) EnableWarn(enable bool) {
	lg.disableWarn = !enable
}

func (lg *logger) Warn(args ...any) {
	if lg.disableWarn {
		return
	}
	lg.print(loggerDepth, &warnLevel, args...)
}

func (lg *logger) Warnf(format string, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.printf(loggerDepth, &warnLevel, &format, args...)
}

func (lg *logger) WarnDepth(depth int, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.print(loggerDepth+depth, &warnLevel, args...)
}

func (lg *logger) WarnfDepth(depth int, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.printf(loggerDepth+depth, &warnLevel, &format, args...)
}

func (lg *logger) WarnTrace(traceID string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &warnLevel, args...)
	} else {
		lg.print(loggerDepth, &warnLevel, args...)
	}
}

func (lg *logger) WarnfTrace(traceID, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &warnLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &warnLevel, &format, args...)
	}
}

func (lg *logger) WarnDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &warnLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &warnLevel, args...)
	}
}

func (lg *logger) WarnfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &warnLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &warnLevel, &format, args...)
	}
}

func (lg *logger) IsError() bool {
	return !lg.disableError
}

func (lg *logger) EnableError(enable bool) {
	lg.disableError = !enable
}

func (lg *logger) Error(args ...any) {
	if lg.disableError {
		return
	}
	lg.print(loggerDepth, &errorLevel, args...)
}

func (lg *logger) Errorf(format string, args ...any) {
	if lg.disableError {
		return
	}
	lg.printf(loggerDepth, &errorLevel, &format, args...)
}

func (lg *logger) ErrorDepth(depth int, args ...any) {
	if lg.disableError {
		return
	}
	lg.print(loggerDepth+depth, &errorLevel, args...)
}

func (lg *logger) ErrorfDepth(depth int, format string, args ...any) {
	if lg.disableError {
		return
	}
	lg.printf(loggerDepth+depth, &errorLevel, &format, args...)
}

func (lg *logger) ErrorTrace(traceID string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &errorLevel, args...)
	} else {
		lg.print(loggerDepth, &errorLevel, args...)
	}
}

func (lg *logger) ErrorfTrace(traceID, format string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &errorLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &errorLevel, &format, args...)
	}
}

func (lg *logger) ErrorDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &errorLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &errorLevel, args...)
	}
}

func (lg *logger) ErrorfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &errorLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &errorLevel, &format, args...)
	}
}
