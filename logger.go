package log

import (
	"fmt"
	"io"
	"runtime"
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
	SetDebug(enable bool)
	Debug(args ...any)
	Debugf(format string, args ...any)
	DebugDepth(depth int, args ...any)
	DebugfDepth(depth int, format string, args ...any)
	DebugTrace(traceID string, args ...any)
	DebugfTrace(traceID string, format string, args ...any)
	DebugDepthTrace(depth int, traceID string, args ...any)
	DebugfDepthTrace(depth int, traceID string, format string, args ...any)
	// // Info 级别的方法。
	// IsInfo() bool
	// Info(args ...any)
	// Infof(format string, args ...any)
	// InfoDepth(depth int, args ...any)
	// InfofDepth(depth int, format string, args ...any)
	// InfoTrace(traceID string, args ...any)
	// InfofTrace(traceID string, format string, args ...any)
	// InfoDepthTrace(traceID string, depth int, args ...any)
	// InfofDepthTrace(traceID string, depth int, format string, args ...any)
	// // Warn 级别的方法。
	// IsWarn() bool
	// Warn(args ...any)
	// Warnf(format string, args ...any)
	// WarnDepth(depth int, args ...any)
	// WarnfDepth(depth int, format string, args ...any)
	// WarnTrace(traceID string, args ...any)
	// WarnfTrace(traceID string, format string, args ...any)
	// WarnDepthTrace(traceID string, depth int, args ...any)
	// WarnfDepthTrace(traceID string, depth int, format string, args ...any)
	// // Error 级别的方法。
	// IsError() bool
	// Error(args ...any)
	// Errorf(format string, args ...any)
	// ErrorDepth(depth int, args ...any)
	// ErrorfDepth(depth int, format string, args ...any)
	// ErrorTrace(traceID string, args ...any)
	// ErrorfTrace(traceID string, format string, args ...any)
	// ErrorDepthTrace(traceID string, depth int, args ...any)
	// ErrorfDepthTrace(traceID string, depth int, format string, args ...any)
}

// logger 默认实现
type logger struct {
	// 输出
	io.Writer
	// 头格式
	HeaderFunc
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
func NewLogger(writer io.Writer, headerFunc HeaderFunc, name string) Logger {
	lg := new(logger)
	lg.Writer = writer
	lg.HeaderFunc = headerFunc
	if name != "" {
		lg.name = fmt.Sprintf("[%s] ", name)
	}
	return lg
}

func (lg *logger) print(depth int, level *string, args ...any) {
	buf := bufPool.Get().(*Buffer)
	buf.b = buf.b[:0]
	// 名称
	if lg.name != "" {
		buf.b = append(buf.b, lg.name...)
	}
	// 级别
	buf.b = append(buf.b, *level...)
	// 头
	lg.HeaderFunc(buf, depth)
	buf.b = append(buf.b, ' ')
	// 日志
	fmt.Fprint(buf, args...)
	// 换行
	buf.b = append(buf.b, '\n')
	// 输出
	lg.Writer.Write(buf.b)
	// 回收
	bufPool.Put(buf)
}

func (lg *logger) printf(depth int, level, format *string, args ...any) {
	buf := bufPool.Get().(*Buffer)
	buf.b = buf.b[:0]
	// 名称
	if lg.name != "" {
		buf.b = append(buf.b, lg.name...)
	}
	// 级别
	buf.b = append(buf.b, *level...)
	// 头
	lg.HeaderFunc(buf, depth)
	buf.b = append(buf.b, ' ')
	// 日志
	fmt.Fprintf(buf, *format, args...)
	// 换行
	buf.b = append(buf.b, '\n')
	// 输出
	lg.Writer.Write(buf.b)
	// 回收
	bufPool.Put(buf)
}

func (lg *logger) printTrace(depth int, trace, level *string, args ...any) {
	buf := bufPool.Get().(*Buffer)
	buf.b = buf.b[:0]
	// 名称
	if lg.name != "" {
		buf.b = append(buf.b, lg.name...)
	}
	// 级别
	buf.b = append(buf.b, *level...)
	// 头
	lg.HeaderFunc(buf, depth)
	buf.b = append(buf.b, ' ')
	// 追踪
	buf.b = append(buf.b, '[')
	buf.b = append(buf.b, *trace...)
	buf.b = append(buf.b, ']')
	buf.b = append(buf.b, ' ')
	// 日志
	fmt.Fprint(buf, args...)
	// 换行
	buf.b = append(buf.b, '\n')
	// 输出
	lg.Writer.Write(buf.b)
	// 回收
	bufPool.Put(buf)
}

func (lg *logger) printfTrace(depth int, trace, level, format *string, args ...any) {
	buf := bufPool.Get().(*Buffer)
	buf.b = buf.b[:0]
	// 名称
	if lg.name != "" {
		buf.b = append(buf.b, lg.name...)
	}
	// 级别
	buf.b = append(buf.b, *level...)
	// 头
	lg.HeaderFunc(buf, depth)
	buf.b = append(buf.b, ' ')
	// 追踪
	buf.b = append(buf.b, '[')
	buf.b = append(buf.b, *trace...)
	buf.b = append(buf.b, ']')
	buf.b = append(buf.b, ' ')
	// 日志
	fmt.Fprintf(buf, *format, args...)
	// 换行
	buf.b = append(buf.b, '\n')
	// 输出
	lg.Writer.Write(buf.b)
	// 回收
	bufPool.Put(buf)
}

func (lg *logger) Recover(recover any) {
	if recover == nil {
		return
	}
	// get statck info buf.line
	buf := bufPool.Get().(*Buffer)
	buf.b = buf.b[:cap(buf.b)]
	for {
		n := runtime.Stack(buf.b, false)
		if n < len(buf.b) {
			buf.b = buf.b[:n]
			break
		}
		buf.b = make([]byte, len(buf.b)+1024)
	}
	// log
	_log := bufPool.Get().(*Buffer)
	_log.b = _log.b[:0]
	// lg.Header.FormatWith(_log, "", PanicLevel, string(b[:i]), string(b[i+1:])
	// fmt.Fprintf(_log, "%v")
	// _log.line = append(_log.line, fmt.sp...)
	// filter
	b := buf.b
	for len(b) > 0 {
		// find new line
		i := 0
		for ; i < len(b); i++ {
			if b[i] == '\n' {
				i++
				break
			}
		}
		line := b[:i]
		b = b[i:]
		// find file line
		if line[0] != '\t' {
			continue
		}
		// \t filepath/file.go:line +0x622
		i = len(line)
		for i > 0 {
			if line[i] == ' ' {
				//
				line = line[:i]
				break
			}
			i--
		}
		// write
		_log.b = append(_log.b, line[1:]...)
	}
	// b := debug.Stack()
	// found := false
	// n := 0
	// for len(b) > 0 {
	// 	i := bytes.IndexByte(b, '\n')
	// 	if i < 0 {
	// 		return
	// 	}
	// 	if !found {
	// 		found = bytes.Contains(b[:i], runtimePanic)
	// 	} else {
	// 		n++
	// 		// the second line
	// 		if n == 2 {
	// 			// skip '\t'
	// 			b := b[1:i]
	// 			_log := logPool.Get().(*Log)
	// 			_log.line = _log.line[:0]
	// 			// stack
	// 			i = bytes.LastIndexByte(b, ' ')
	// 			if i > 0 {
	// 				b = b[:i]
	// 				// split path and line
	// 				i = bytes.LastIndexByte(b, ':')
	// 				if i > 0 {
	// 					lg.Header.FormatWith(_log, "", PanicLevel, string(b[:i]), string(b[i+1:]))
	// 				} else {
	// 					lg.Header.FormatWith(_log, "", PanicLevel, "?", "-1")
	// 				}
	// 			} else {
	// 				lg.Header.FormatWith(_log, "", PanicLevel, "?", "-1")
	// 			}
	// 			//
	// 			_log.line = append(_log.line, ' ')
	// 			fmt.Fprint(_log, recover)
	// 			_log.line = append(_log.line, '\n')
	// 			lg.Writer.Write(_log.line)
	// 			logPool.Put(_log)
	// 			return
	// 		}
	// 	}
	// 	b = b[i+1:]
	// }
}

func (lg *logger) IsDebug() bool {
	return !lg.disableDebug
}

func (lg *logger) SetDebug(enable bool) {
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
