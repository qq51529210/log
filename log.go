package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	logPool                     = sync.Pool{} // Logger缓存
	SpaceSeparator              = " "         // 空格
	DateSeparator               = "-"         // 日期
	TimeSeparator               = ":"         // 时间
	NanoSecSeparator            = "."         // 纳秒
	FileLineSeparator           = ":"         // 堆栈
	NanosecondLength            = 6           // 纳秒的格式化长度
	defaultWriter     io.Writer = os.Stdout   // 默认输出
	DebugLevel                  = "D"
	InfoLevel                   = "I"
	WarnLevel                   = "W"
	ErrorLevel                  = "E"
	PanicLevel                  = "P"
	endLine                     = []byte("\n")
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
	b []byte // 缓存
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
	// 值是倒转的
	i1 := len(l.b)
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
		i1++
	}
	end := i1 + length
	for {
		l.b = append(l.b, byte('0'+integer%10))
		integer /= 10
		if integer == 0 {
			break
		}
		length--
	}
	// 先反转
	c := byte(0)
	i2 := len(l.b) - 1
	for i1 < i2 {
		c = l.b[i1]
		l.b[i1] = l.b[i2]
		l.b[i2] = c
		i2--
		i1++
	}
	// 再后面补0
	for length > 0 {
		l.b = append(l.b, byte('0'))
		length--
	}
	if len(l.b) > end {
		l.b = l.b[:end]
	}
}

// 写入integer，长度是length，不够左边补0
func (l *Log) IntL0(integer, length int) {
	// 值是倒转的
	i1 := len(l.b)
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
		i1++
	}
	for {
		l.b = append(l.b, byte('0'+integer%10))
		integer /= 10
		if integer == 0 {
			break
		}
		length--
	}
	// 继续在后面补0
	for length > 1 {
		l.b = append(l.b, byte('0'))
		length--
	}
	// 反转
	i2 := len(l.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = l.b[i1]
		l.b[i1] = l.b[i2]
		l.b[i2] = c
		i2--
		i1++
	}
}

// 写入一个整数integer
func (l *Log) Int(integer int) {
	i1 := len(l.b)
	if integer < 0 {
		l.b = append(l.b, '-')
		integer = -integer
		i1++
	}
	for {
		l.b = append(l.b, byte('0'+integer%10))
		integer /= 10
		if integer == 0 {
			break
		}
	}
	i2 := len(l.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = l.b[i1]
		l.b[i1] = l.b[i2]
		l.b[i2] = c
		i2--
		i1++
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
	l.IntL0(year, 4)
	l.b = append(l.b, DateSeparator...)
	l.IntL0(int(month), 2)
	l.b = append(l.b, DateSeparator...)
	l.IntL0(day, 2)
	l.b = append(l.b, SpaceSeparator...)
	l.IntL0(hour, 2)
	l.b = append(l.b, TimeSeparator...)
	l.IntL0(minute, 2)
	l.b = append(l.b, TimeSeparator...)
	l.IntL0(second, 2)
	if NanosecondLength > 0 {
		l.b = append(l.b, NanoSecSeparator...)
		l.IntR0(t.Nanosecond(), NanosecondLength)
	}
}

// 写入堆栈信息
func (l *Log) Stack() {
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
}

// level:日志级别，skip:堆栈调用层级
func (l *Log) Header(level string, skip int) {
	// 级别
	l.b = append(l.b, level...)
	l.b = append(l.b, SpaceSeparator...)
	// 时间
	l.Time()
	l.b = append(l.b, SpaceSeparator...)
	// 调用堆栈
	l.PathLine(skip + 1)
	l.b = append(l.b, SpaceSeparator...)
}

// 写入堆栈的文件的完整路径path的文件名和行号line，skip:堆栈的层次
func (l *Log) FileLine(skip int) {
	_, path, line, o := runtime.Caller(skip + 1)
	if !o {
		l.b = append(l.b, "???"...)
		l.b = append(l.b, FileLineSeparator...)
		l.Int(-1)
	} else {
		i := len(path) - 1
		for ; i >= 0; i-- {
			if os.IsPathSeparator(path[i]) {
				i++
				break
			}
		}
		l.b = append(l.b, path[i:]...)
		l.b = append(l.b, FileLineSeparator...)
		l.Int(line)
	}
}

// 写入堆栈的文件的完整路径path和行号line，skip:堆栈的层次
func (l *Log) PathLine(skip int) {
	_, path, line, o := runtime.Caller(skip + 1)
	if !o {
		l.b = append(l.b, "???"...)
		l.b = append(l.b, FileLineSeparator...)
		l.Int(-1)
	} else {
		l.b = append(l.b, path...)
		l.b = append(l.b, FileLineSeparator...)
		l.Int(line)
	}
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
	l.Header(level, skip+1)
	l.b = append(l.b, str...)
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，data:数据
func PrintBytes(level string, skip int, data []byte) {
	l := GetLog()
	l.Header(level, skip+1)
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	_, _ = defaultWriter.Write(data)
	_, _ = defaultWriter.Write(endLine)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，format:格式化字符串，args:格式化数据
func Printf(level string, skip int, format string, args ...interface{}) {
	l := GetLog()
	l.Header(level, skip+1)
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	_, _ = fmt.Fprintf(defaultWriter, format, args...)
	_, _ = defaultWriter.Write(endLine)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，args:格式化数据
func Fprint(level string, skip int, args ...interface{}) {
	l := GetLog()
	l.Header(level, skip+1)
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	_, _ = fmt.Fprint(defaultWriter, args...)
	_, _ = defaultWriter.Write(endLine)
}

// 使用defaultWriter
func Debug(a ...interface{}) {
	Fprint(DebugLevel, 1, a...)
}

// 使用defaultWriter
func Info(a ...interface{}) {
	Fprint(InfoLevel, 1, a...)
}

// 使用defaultWriter
func Warn(a ...interface{}) {
	Fprint(WarnLevel, 1, a...)
}

// 使用defaultWriter
func Error(a ...interface{}) {
	Fprint(ErrorLevel, 1, a...)
}

// 如果recover调用函数f
func Recover(f func()) {
	re := recover()
	if re == nil {
		return
	}
	l := GetLog()
	l.b = append(l.b, PanicLevel...)
	l.b = append(l.b, SpaceSeparator...)
	l.Time()
	l.b = append(l.b, SpaceSeparator...)
	l.Stack()
	l.EndLine()
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	if f != nil {
		f()
	}
}
