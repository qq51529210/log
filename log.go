package log

import "sync"

var (
	logPool = sync.Pool{}
)

func init() {
	logPool.New = func() interface{} {
		return new(Log)
	}
}

type Error string

func (e Error) Error() string {
	return string(e)
}

// Get a Log from pool.
func GetLog() *Log {
	return logPool.Get().(*Log)
}

// Put Log into pool.
func PutLog(l *Log) {
	l.Reset()
	logPool.Put(l)
}

// A buffer used for format a line of logs.
type Log struct {
	// Log buffer.
	line []byte
	// Integer format buffer.
	buff []byte
}

// Implements io.Writer interface.
func (l *Log) Write(b []byte) (int, error) {
	l.line = append(l.line, b...)
	return len(b), nil
}

// Return buffer.
func (l *Log) Data() []byte {
	return l.line
}

// Return string of buffer.
func (l *Log) String() string {
	return string(l.line)
}

// Reset buffer.
func (l *Log) Reset() {
	l.line = l.line[:0]
}

// Write a integer into buffer with right align format.
// If len(integer) < length,add 0 to the left.
// Example: 12 -> 0012 while length=4.
func (l *Log) WriteRightAlignInt(integer, length int) {
	// Zero.
	if integer == 0 {
		for i := 0; i < length; i++ {
			l.line = append(l.line, '0')
		}
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.buff = l.buff[:0]
	for integer > 0 {
		l.buff = append(l.buff, byte('0'+integer%10))
		integer /= 10
	}
	// Add 0 to the left if len(integer) < length
	if length > len(l.buff) {
		for i := len(l.buff); i < length; i++ {
			l.line = append(l.line, '0')
		}
	}
	// buff[4,3,2,1]->line[1,2,3,4]
	for i := len(l.buff) - 1; i >= 0; i-- {
		l.line = append(l.line, l.buff[i])
	}
}

// Write a integer into buffer with left align format.
// If len(integer) < length,add 0 to the right.
// Example:
//	12 -> 1200 while length=4.
//	1234 -> 12 while length=2.
func (l *Log) WriteLeftAlignInt(integer, length int) {
	// Zero.
	if integer == 0 {
		for i := 0; i < length; i++ {
			l.line = append(l.line, '0')
		}
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.buff = l.buff[:0]
	for integer > 0 {
		l.buff = append(l.buff, byte('0'+integer%10))
		integer /= 10
	}
	if length < len(l.buff) {
		// buff[4,3,2,1]->line[1,2]
		for i := len(l.buff) - 1; i >= len(l.buff)-length; i-- {
			l.line = append(l.line, l.buff[i])
		}
	} else {
		// buff[4,3,2,1]->line[1,2,3,4]
		for i := len(l.buff) - 1; i >= 0; i-- {
			l.line = append(l.line, l.buff[i])
		}
		// Add 0 to the right.
		for i := len(l.buff); i < length; i++ {
			l.line = append(l.line, '0')
		}
	}
}

// Write a integer into buffer without algin format.
func (l *Log) WriteInt(integer int) {
	// Zero.
	if integer == 0 {
		l.line = append(l.line, '0')
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.buff = l.buff[:0]
	for integer > 0 {
		l.buff = append(l.buff, byte('0'+integer%10))
		integer /= 10
	}
	// buff[4,3,2,1]->line[1,2,3,4]
	for i := len(l.buff) - 1; i >= 0; i-- {
		l.line = append(l.line, l.buff[i])
	}
}

// Write a byte into buffer.
func (l *Log) WriteUint8(c byte) {
	l.line = append(l.line, c)
}

// Write binary array into buffer.
func (l *Log) WriteBytes(b []byte) {
	l.line = append(l.line, b...)
}

// Write a string into buffer.
func (l *Log) WriteString(s string) {
	l.line = append(l.line, s...)
}

// // 写入堆栈信息
// func (l *Log) Stack() {
// 	// 所有的堆栈
// 	i1 := len(l.b)
// 	i2 := 0
// 	n := 0
// 	for {
// 		l.b = l.b[:cap(l.b)]
// 		n = runtime.Stack(l.b[i1:], true)
// 		i2 = i1 + n
// 		if i2 < len(l.b) {
// 			l.b = l.b[:i2]
// 			break
// 		}
// 		l.b = append(l.b, make([]byte, 128)...)
// 	}
// 	/*
// 		github.com/qq51529210/log.(*Log).Stack(0xc00000c080)
// 		        /Users/ben/Documents/project/go/src/github.com/qq51529210/log/log.go:221 +0x7f
// 		github.com/qq51529210/log.Recover(0x10d5c70)
// 		        /Users/ben/Documents/project/go/src/github.com/qq51529210/log/log.go:357 +0x1fa
// 		panic(0x10b4540, 0x10ed8e0)
// 		        /usr/local/go/src/runtime/panic.go:969 +0x175
// 		main.f1(...)
// 		        /Users/ben/Documents/project/go/src/test/main.go:20
// 		main.main()
// 		        /Users/ben/Documents/project/go/src/test/main.go:27 +0x65
// 	*/
// 	// 简化一下，只保留文件路径
// 	n = i2 - 1
// 	i := i1
// 	m := i
// 	// 是否找到/runtime/panic.go，下一行就是panic的地方
// 	ok := false
// 	for i < n {
// 		// 文件行开始，'\t'
// 		if l.b[i] == '\t' {
// 			i++
// 			i1 = i
// 			for ; i < n; i++ {
// 				// 文件行路径结束
// 				if l.b[i] == ' ' || l.b[i] == '\n' {
// 					// 找到/runtime/panic.go
// 					if !ok {
// 						i2 = i - 1
// 						i2 = i1 + bytes.LastIndexByte(l.b[i1:i2], ':')
// 						if i2 > 0 {
// 							ok = bytes.LastIndex(l.b[i1:i2], panicFileLine) >= 0
// 						}
// 					} else {
// 						m += copy(l.b[m:], l.b[i1:i])
// 						l.b[m] = ' '
// 						m++
// 					}
// 					break
// 				}
// 			}
// 		}
// 		i++
// 	}
// 	l.b = l.b[:m]
// 	return
// }

// // 写入堆栈的文件的完整路径path的文件名和行号line，skip:堆栈的层次
// func (l *Log) FileLine(path string, line int) {
// 	i := len(path) - 1
// 	for ; i >= 0; i-- {
// 		if os.IsPathSeparator(path[i]) {
// 			i++
// 			break
// 		}
// 	}
// 	l.b = append(l.b, path[i:]...)
// 	l.b = append(l.b, FileLineSeparator)
// 	l.Int(line)
// }

// // 写入堆栈的文件的完整路径path和行号line，skip:堆栈的层次
// func (l *Log) PathLine(path string, line int) {
// 	l.b = append(l.b, path...)
// 	l.b = append(l.b, FileLineSeparator)
// 	l.Int(line)
// }

// // 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，str:日志文本
// func Print(level string, skip int, str string) {
// 	_, path, line, ok := runtime.Caller(skip + 1)
// 	if !ok {
// 		path = "???"
// 		line = -1
// 	}
// 	l := GetLog()
// 	defaultPrintHeader(l, level, path, line)
// 	l.b = append(l.b, str...)
// 	l.EndLine()
// 	_, _ = defaultWriter.Write(l.b)
// 	logPool.Put(l)
// }

// // 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，data:数据
// func PrintBytes(level string, skip int, data []byte) {
// 	_, path, line, ok := runtime.Caller(skip + 1)
// 	if !ok {
// 		path = "???"
// 		line = -1
// 	}
// 	l := GetLog()
// 	defaultPrintHeader(l, level, path, line)
// 	l.b = append(l.b, data...)
// 	l.EndLine()
// 	_, _ = defaultWriter.Write(l.b)
// 	logPool.Put(l)
// }

// // 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，format:格式化字符串，args:格式化数据
// func Printf(level string, skip int, format string, args ...interface{}) {
// 	_, path, line, ok := runtime.Caller(skip + 1)
// 	if !ok {
// 		path = "???"
// 		line = -1
// 	}
// 	l := GetLog()
// 	defaultPrintHeader(l, level, path, line)
// 	_, _ = fmt.Fprintf(l, format, args...)
// 	l.EndLine()
// 	_, _ = defaultWriter.Write(l.b)
// 	logPool.Put(l)
// }

// // 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，args:格式化数据
// func Fprint(level string, skip int, args ...interface{}) {
// 	_, path, line, ok := runtime.Caller(skip + 1)
// 	if !ok {
// 		path = "???"
// 		line = -1
// 	}
// 	l := GetLog()
// 	defaultPrintHeader(l, level, path, line)
// 	_, _ = fmt.Fprint(l, args...)
// 	l.EndLine()
// 	_, _ = defaultWriter.Write(l.b)
// 	logPool.Put(l)
// }

// func Debug(a ...interface{}) {
// 	Fprint(DebugLevel, 1, a...)
// }

// func Info(a ...interface{}) {
// 	Fprint(InfoLevel, 1, a...)
// }

// func Warn(a ...interface{}) {
// 	Fprint(WarnLevel, 1, a...)
// }

// func Error(a ...interface{}) {
// 	Fprint(ErrorLevel, 1, a...)
// }

// func DebugSkip(skip int, a ...interface{}) {
// 	Fprint(DebugLevel, skip+1, a...)
// }

// func InfoSkip(skip int, a ...interface{}) {
// 	Fprint(InfoLevel, skip+1, a...)
// }

// func WarnSkip(skip int, a ...interface{}) {
// 	Fprint(WarnLevel, skip+1, a...)
// }

// func ErrorSkip(skip int, a ...interface{}) {
// 	Fprint(ErrorLevel, skip+1, a...)
// }

// // 如果recover调用函数f
// func Recover(f func()) {
// 	re := recover()
// 	if re == nil {
// 		return
// 	}
// 	l := GetLog()
// 	l.b = append(l.b, PanicLevel...)
// 	l.b = append(l.b, SpaceSeparator)
// 	l.Time()
// 	l.b = append(l.b, SpaceSeparator)
// 	l.Stack()
// 	l.b = append(l.b, SpaceSeparator)
// 	_, _ = fmt.Fprint(l, re)
// 	l.EndLine()
// 	_, _ = defaultWriter.Write(l.b)
// 	logPool.Put(l)
// 	if f != nil {
// 		f()
// 	}
// }
