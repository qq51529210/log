package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level byte

const (
	LevelDebug Level = 'D'
	LevelInfo  Level = 'I'
	LevelWarn  Level = 'W'
	LevelError Level = 'E'
	LevelPanic Level = 'P'
)

////type StackInfo byte
////
////const (
////	StackInfoDisable StackInfo = iota // 不打印堆栈信息
////	StackInfoFile                     // 打印文件名称
////	StackInfoPath                     // 打印文件的完整路径
////)
//

var (
	logPool                     = sync.Pool{} // Logger缓存
	SpaceSeparator    byte      = ' '         // 空格
	DateSeparator     byte      = '-'         // 日期
	TimeSeparator     byte      = ':'         // 时间
	NanoSecSeparator  byte      = '.'         // 纳秒
	FileLineSeparator byte      = ':'         // 堆栈
	defaultWriter     io.Writer = os.Stdout   // 默认输出
	////	//panicFileLine               = []byte("/runtime/panic.go") // 获取panic堆栈判断（只获取panic那行）
)

// 缓存池的new函数
func init() {
	logPool.New = func() interface{} {
		return &Log{}
	}
}

func GetLog() *Log {
	return logPool.Get().(*Log)
}

func PutLog(l *Log) {
	logPool.Put(l)
}

// 一行日志数据缓存
type Log struct {
	b []byte // 缓存
}

// io.WriteTo接口
func (l *Log) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(l.b)
	if err != nil {
		return 0, err
	}
	l.Reset()
	return int64(n), err
}

// io.Writer接口
func (l *Log) Write(b []byte) (int, error) {
	l.b = append(l.b, b...)
	return len(b), nil
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

// 写入换行'\n'
func (l *Log) EndLine() {
	l.b = append(l.b, '\n')
}

// 写入空格' '
func (l *Log) Space() {
	l.b = append(l.b, SpaceSeparator)
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
func (l *Log) Time(dateTime *time.Time) {
	year, month, day := dateTime.Date()
	hour, minute, second := dateTime.Clock()
	l.IntL0(year, 4)
	l.b = append(l.b, DateSeparator)
	l.IntL0(int(month), 2)
	l.b = append(l.b, DateSeparator)
	l.IntL0(day, 2)
	l.b = append(l.b, SpaceSeparator)
	l.IntL0(hour, 2)
	l.b = append(l.b, TimeSeparator)
	l.IntL0(minute, 2)
	l.b = append(l.b, TimeSeparator)
	l.IntL0(second, 2)
	l.b = append(l.b, NanoSecSeparator)
	l.IntR0(dateTime.Nanosecond(), 9)
}

// 写入日志级别level
func (l *Log) Level(level Level) {
	l.b = append(l.b, byte(level))
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
func (l *Log) Header(level Level, skip int) {
	l.Reset()
	// 级别
	l.Level(level)
	l.Space()
	// 日期，时间
	t := time.Now()
	l.Time(&t)
	l.Space()
	//
	l.PathLine(skip + 1)
	l.Space()
}

// 写入堆栈的文件的完整路径path的文件名和行号line，skip:堆栈的层次
func (l *Log) FileLine(skip int) {
	_, path, line, o := runtime.Caller(skip + 1)
	if !o {
		l.b = append(l.b, "???"...)
		l.b = append(l.b, FileLineSeparator)
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
		l.b = append(l.b, FileLineSeparator)
		l.Int(line)
	}
}

// 写入堆栈的文件的完整路径path和行号line，skip:堆栈的层次
func (l *Log) PathLine(skip int) {
	_, path, line, o := runtime.Caller(skip)
	if !o {
		l.b = append(l.b, "???"...)
		l.b = append(l.b, FileLineSeparator)
		l.Int(-1)
	} else {
		l.b = append(l.b, path...)
		l.b = append(l.b, FileLineSeparator)
		l.Int(line)
	}
}

// w:输出，level:日志级别，skip:堆栈调用层级，str:日志文本
func (l *Log) Print(w io.Writer, level Level, skip int, str string) (int, error) {
	l.Header(level, skip+1)
	l.String(str)
	l.EndLine()
	return w.Write(l.b)
}

// w:输出，level:日志级别，skip:堆栈调用层级，data:二进制数据
func (l *Log) PrintBytes(w io.Writer, level Level, skip int, data []byte) (int, error) {
	l.Header(level, skip+1)
	l.Bytes(data)
	l.EndLine()
	return w.Write(l.b)
}

// w:输出，level:日志级别，skip:堆栈调用层级，format:格式化字符串，args:格式化数据
func (l *Log) Printf(w io.Writer, level Level, skip int, format string, args ...interface{}) (int, error) {
	l.Header(level, skip+1)
	_, _ = fmt.Fprintf(l, format, args...)
	l.EndLine()
	return w.Write(l.b)
}

// w:输出，level:日志级别，skip:堆栈调用层级，args:格式化数据
func (l *Log) Fprint(w io.Writer, level Level, skip int, args ...interface{}) (int, error) {
	l.Header(level, skip+1)
	_, _ = fmt.Fprint(l, args...)
	l.EndLine()
	return w.Write(l.b)
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，str:日志文本
func Print(level Level, skip int, str string) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Print(defaultWriter, level, skip+1, str)
	logPool.Put(l)
	return n, e
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，format:格式化字符串，args:格式化数据
func Printf(level Level, skip int, format string, args ...interface{}) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Printf(defaultWriter, level, skip+1, format, args...)
	logPool.Put(l)
	return n, e
}

// 使用默认writer输出日志，level:日志级别，skip:堆栈调用层级，args:格式化数据
func Fprint(level Level, skip int, args ...interface{}) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Fprint(defaultWriter, level, skip+1, args...)
	logPool.Put(l)
	return n, e
}

// 使用defaultWriter和defaultStack
func Debug(a ...interface{}) {
	_, _ = Fprint(LevelDebug, 1, a...)
}

// 使用defaultWriter和defaultStack
func Info(a ...interface{}) {
	_, _ = Fprint(LevelInfo, 1, a...)
}

// 使用defaultWriter和defaultStack
func Warn(a ...interface{}) {
	_, _ = Fprint(LevelWarn, 1, a...)
}

// 使用defaultWriter和defaultStack
func Error(a ...interface{}) {
	_, _ = Fprint(LevelError, 1, a...)
}

// 设置默认的输出writer
func SetWriter(w io.Writer) {
	if w != nil {
		defaultWriter = w
	}
}

// 如果recover调用函数f
func Recover(f func()) {
	re := recover()
	if re == nil {
		return
	}
	//
	l := logPool.Get().(*Log)
	l.Reset()
	l.Level(LevelPanic)
	l.Space()
	t := time.Now()
	l.Time(&t)
	l.Space()
	l.Stack()
	//
	_, _ = defaultWriter.Write(l.b)
	logPool.Put(l)
	//
	if f != nil {
		f()
	}
}
