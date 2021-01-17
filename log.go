package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	logPool                 = sync.Pool{} // Log缓存池
	defaultWriter io.Writer = os.Stdout   // 默认输出
	//SpaceSeparator              = " "         // 空格
	//DateSeparator               = "-"         // 日期
	//DateTimeSpace               = " "         // 日期和时间之间的空格
	//TimeSeparator               = ":"         // 时间
	//NanoSecSeparator            = "."         // 纳秒
	//FileLineSeparator           = ":"         // 堆栈
	SpaceSeparator    byte = ' '                        // 空格
	DateSeparator     byte = '-'                        // 日期
	DateTimeSpace     byte = ' '                        // 日期和时间之间的空格
	TimeSeparator     byte = ':'                        // 时间
	NanoSecSeparator  byte = '.'                        // 纳秒
	FileLineSeparator byte = ':'                        // 堆栈
	NanosecondLength       = 6                          // 纳秒的格式化长度
	DebugLevel             = "D"                        // Debug函数的级别
	InfoLevel              = "I"                        // Info函数的级别
	WarnLevel              = "W"                        // Warn函数的级别
	ErrorLevel             = "E"                        // Error函数的级别
	PanicLevel             = "P"                        // Recover函数的级别
	panicFileLine          = []byte("runtime/panic.go") // Stack函数查找panic的标志
)

const (
	BaseSkip = 1 // Print/Printf/Fprint函数内有1层调用
)

// 缓存池的new函数
func init() {
	logPool.New = func() interface{} {
		return new(Log)
	}
}

// 获取日志缓存
func GetLog() *Log {
	l := logPool.Get().(*Log)
	l.Reset()
	return l
}

// 回收日志缓存
func PutLog(l *Log) {
	logPool.Put(l)
}

// 一行日志数据缓存
type Log struct {
	b []byte   // 缓存
	i [20]byte // 整数格式化缓存
}

// io.Writer接口
func (l *Log) Write(b []byte) (int, error) {
	l.b = append(l.b, b...)
	return len(b), nil
}

// 返回数据缓存
func (l *Log) Data() []byte {
	return l.b
}

// 重置缓存
func (l *Log) Reset() {
	l.b = l.b[:0]
}

// 写入integer，长度是length，不够右边补0
func (l *Log) IntR0(integer, length int) {
	i := len(l.i) - 1
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
	}
	for {
		l.i[i] = byte('0' + integer%10)
		integer /= 10
		length--
		if integer < 1 {
			l.b = append(l.b, l.i[i:]...)
			for length > 0 {
				l.b = append(l.b, '0')
				length--
			}
			return
		}
		i--
	}
}

// 写入integer，长度是length，不够左边补0
func (l *Log) IntL0(integer, length int) {
	i := len(l.i) - 1
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
	}
	for {
		l.i[i] = byte('0' + integer%10)
		integer /= 10
		length--
		if integer < 1 {
			for length > 0 {
				length--
				i--
				l.i[i] = '0'
			}
			l.b = append(l.b, l.i[i:]...)
			return
		}
		i--
	}
}

// 写入一个整数integer
func (l *Log) Int(integer int) {
	i := len(l.i) - 1
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
	}
	for {
		l.i[i] = byte('0' + integer%10)
		integer /= 10
		if integer < 1 {
			l.b = append(l.b, l.i[i:]...)
			return
		}
		i--
	}
}

// 写入一个字符c
func (l *Log) Byte(c byte) {
	l.b = append(l.b, c)
}

// 写入空白字符
func (l *Log) Space() {
	l.b = append(l.b, ' ')
}

// 写入换行'\n'
func (l *Log) EndLine() {
	l.b = append(l.b, '\n')
}

// 写入字符串str
func (l *Log) String(str string) {
	l.b = append(l.b, str...)
}

// 写入二进制数据data
func (l *Log) Bytes(data []byte) {
	l.b = append(l.b, data...)
}

// 写入日期时间dateTime，格式（year-month-day hour:minute:second.nano），如果NanoSecLength=0则不写入纳秒
func (l *Log) Time() {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// 日期
	l.IntL0(year, 4)
	//l.b = append(l.b, DateSeparator...)
	l.b = append(l.b, DateSeparator)
	l.IntL0(int(month), 2)
	//l.b = append(l.b, DateSeparator...)
	l.b = append(l.b, DateSeparator)
	l.IntL0(day, 2)
	//l.b = append(l.b, DateTimeSpace...)
	l.b = append(l.b, DateTimeSpace)
	// 时间
	l.IntL0(hour, 2)
	//l.b = append(l.b, TimeSeparator...)
	l.b = append(l.b, TimeSeparator)
	l.IntL0(minute, 2)
	//l.b = append(l.b, TimeSeparator...)
	l.b = append(l.b, TimeSeparator)
	l.IntL0(second, 2)
	// 纳秒
	if NanosecondLength > 0 {
		//l.b = append(l.b, NanoSecSeparator...)
		l.b = append(l.b, NanoSecSeparator)
		l.IntR0(t.Nanosecond(), NanosecondLength)
	}
}

// 写入堆栈信息
func (l *Log) Stack() {
	// 所有的堆栈
	i1 := len(l.b)
	i2 := 0
	n := 0
	for {
		l.b = l.b[:cap(l.b)]
		n = runtime.Stack(l.b[i1:], true)
		i2 = i1 + n
		if i2 < len(l.b) {
			l.b = l.b[:i2]
			break
		}
		l.b = append(l.b, make([]byte, 128)...)
	}
	/*
		github.com/qq51529210/log.(*Log).Stack(0xc00000c080)
		        /Users/ben/Documents/project/go/src/github.com/qq51529210/log/log.go:221 +0x7f
		github.com/qq51529210/log.Recover(0x10d5c70)
		        /Users/ben/Documents/project/go/src/github.com/qq51529210/log/log.go:357 +0x1fa
		panic(0x10b4540, 0x10ed8e0)
		        /usr/local/go/src/runtime/panic.go:969 +0x175
		main.f1(...)
		        /Users/ben/Documents/project/go/src/test/main.go:20
		main.main()
		        /Users/ben/Documents/project/go/src/test/main.go:27 +0x65
	*/
	// 简化一下，只保留文件路径
	n = i2 - 1
	i := i1
	m := i
	// 是否找到/runtime/panic.go，下一行就是panic的地方
	ok := false
	for i < n {
		// 文件行开始，'\t'
		if l.b[i] == '\t' {
			i++
			i1 = i
			for ; i < n; i++ {
				// 文件行路径结束
				if l.b[i] == ' ' || l.b[i] == '\n' {
					// 找到/runtime/panic.go
					if !ok {
						i2 = i - 1
						i2 = i1 + bytes.LastIndexByte(l.b[i1:i2], ':')
						if i2 > 0 {
							ok = bytes.LastIndex(l.b[i1:i2], panicFileLine) >= 0
						}
					} else {
						m += copy(l.b[m:], l.b[i1:i])
						l.b[m] = ' '
						m++
					}
					break
				}
			}
		}
		i++
	}
	l.b = l.b[:m]
	return
}

// level:日志级别，skip:堆栈调用层级
func (l *Log) Header(level string, path string, line int) {
	// 级别
	l.b = append(l.b, level...)
	//l.b = append(l.b, SpaceSeparator...)
	l.b = append(l.b, SpaceSeparator)
	// 时间
	l.Time()
	//l.b = append(l.b, SpaceSeparator...)
	l.b = append(l.b, SpaceSeparator)
	// 调用堆栈
	l.PathLine(path, line)
	//l.b = append(l.b, SpaceSeparator...)
	l.b = append(l.b, SpaceSeparator)
}

// 写入堆栈的文件的完整路径path的文件名和行号line，skip:堆栈的层次
func (l *Log) FileLine(path string, line int) {
	i := len(path) - 1
	for ; i >= 0; i-- {
		if os.IsPathSeparator(path[i]) {
			i++
			break
		}
	}
	l.b = append(l.b, path[i:]...)
	//l.b = append(l.b, FileLineSeparator...)
	l.b = append(l.b, FileLineSeparator)
	l.Int(line)
}

// 写入堆栈的文件的完整路径path和行号line，skip:堆栈的层次
func (l *Log) PathLine(path string, line int) {
	l.b = append(l.b, path...)
	//l.b = append(l.b, FileLineSeparator...)
	l.b = append(l.b, FileLineSeparator)
	l.Int(line)
}

// 设置默认的输出writer
func SetWriter(w io.Writer) {
	if w != nil {
		defaultWriter = w
	}
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，str:日志文本
func Print(level string, skip int, str string) {
	l := GetLog()
	// 调用堆栈
	_, path, line, o := runtime.Caller(skip)
	if !o {
		path = "???"
		line = -1
	}
	l.Header(level, path, line)
	l.b = append(l.b, str...)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，data:数据
func PrintBytes(level string, skip int, data []byte) {
	l := GetLog()
	// 调用堆栈
	_, path, line, o := runtime.Caller(skip)
	if !o {
		path = "???"
		line = -1
	}
	l.Header(level, path, line)
	l.b = append(l.b, data...)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，format:格式化字符串，args:格式化数据
func Printf(level string, skip int, format string, args ...interface{}) {
	l := GetLog()
	// 调用堆栈
	_, path, line, o := runtime.Caller(skip)
	if !o {
		path = "???"
		line = -1
	}
	l.Header(level, path, line)
	_, _ = fmt.Fprintf(l, format, args...)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，args:格式化数据
func Fprint(level string, skip int, args ...interface{}) {
	l := GetLog()
	// 调用堆栈
	_, path, line, o := runtime.Caller(skip)
	if !o {
		path = "???"
		line = -1
	}
	l.Header(level, path, line)
	_, _ = fmt.Fprint(l, args...)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
}

func Debug(a ...interface{}) {
	Fprint(DebugLevel, 2, a...)
}

func Info(a ...interface{}) {
	Fprint(InfoLevel, 2, a...)
}

func Warn(a ...interface{}) {
	Fprint(WarnLevel, 2, a...)
}

func Error(a ...interface{}) {
	Fprint(ErrorLevel, 2, a...)
}

func DebugSkip(skip int, a ...interface{}) {
	Fprint(DebugLevel, skip+2, a...)
}

func InfoSkip(skip int, a ...interface{}) {
	Fprint(InfoLevel, skip+2, a...)
}

func WarnSkip(skip int, a ...interface{}) {
	Fprint(WarnLevel, skip+2, a...)
}

func ErrorSkip(skip int, a ...interface{}) {
	Fprint(ErrorLevel, skip+2, a...)
}

// 如果recover调用函数f
func Recover(f func()) {
	re := recover()
	if re == nil {
		return
	}
	l := GetLog()
	l.b = append(l.b, PanicLevel...)
	//l.b = append(l.b, SpaceSeparator...)
	l.b = append(l.b, SpaceSeparator)
	l.Time()
	//l.b = append(l.b, SpaceSeparator...)
	l.b = append(l.b, SpaceSeparator)
	l.Stack()
	l.b = append(l.b, SpaceSeparator)
	_, _ = fmt.Fprint(l, re)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	if f != nil {
		f()
	}
}
