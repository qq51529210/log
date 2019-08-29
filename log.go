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

type Level byte

const (
	LevelDebug   Level = 'D'
	LevelInfo    Level = 'I'
	LevelWarn    Level = 'W'
	LevelError   Level = 'E'
	LevelPanic   Level = 'P'
	LevelRecover Level = 'R'
)

type FileLine byte

const (
	FileLineDisable  FileLine = iota
	FileLineFullPath
	FileLineName
)

var (
	logPool                = sync.Pool{}
	errPool                = sync.Pool{}
	SpaceSeparator    byte = ' ' // 空格
	DateSeparator     byte = '-' // 日期
	TimeSeparator     byte = ':' // 时间
	NanoSecSeparator  byte = '.' // 纳秒
	FileLineSeparator byte = ':' // 堆栈
	NanoSecLength          = 6   // 打印纳秒的长度
	unknownFileLine        = []byte("???:-1:")
	panicLine              = []byte("/runtime/panic.go")
)

func init() {
	logPool.New = func() interface{} {
		return &Log{}
	}
	errPool.New = func() interface{} {
		return &panicError{}
	}
}

type Log struct {
	b []byte
}

func (this *Log) Reset() {
	this.b = this.b[0:0]
}

func (this *Log) Bytes() []byte {
	return this.b
}

func (this *Log) IntegerAlignRight(value, length int) {
	// 值是倒转的
	i1 := len(this.b)
	end := i1 + length
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
		length--
	}
	// 先反转
	c := byte(0)
	i2 := len(this.b) - 1
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
	// 再后面补0
	for length > 0 {
		this.b = append(this.b, byte('0'))
		length--
	}
	if len(this.b) > end {
		this.b = this.b[:end]
	}
}

func (this *Log) IntegerAlignLeft(value, length int) {
	// 值是倒转的
	i1 := len(this.b)
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
		length--
	}
	// 继续在后面补0
	i2 := len(this.b) - i1
	for length > i2 {
		this.b = append(this.b, byte('0'))
		length--
	}
	// 反转
	i2 = len(this.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
}

func (this *Log) Integer(value int) {
	i1 := len(this.b)
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
	}
	i2 := len(this.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
}

func (this *Log) Byte(c byte) {
	this.b = append(this.b, c)
}

func (this *Log) EndLine() {
	this.b = append(this.b, '\n')
}

func (this *Log) DateTime(nsec int) {
	date_time := time.Now()
	year, month, day := date_time.Date()
	hour, minute, second := date_time.Clock()

	this.IntegerAlignLeft(year, 4)
	this.b = append(this.b, DateSeparator)
	this.IntegerAlignLeft(int(month), 2)
	this.b = append(this.b, DateSeparator)
	this.IntegerAlignLeft(day, 2)
	this.b = append(this.b, SpaceSeparator)
	this.IntegerAlignLeft(hour, 2)
	this.b = append(this.b, TimeSeparator)
	this.IntegerAlignLeft(minute, 2)
	this.b = append(this.b, TimeSeparator)
	this.IntegerAlignLeft(second, 2)
	if nsec > 0 {
		// 再长没有意义
		if nsec > 6 {
			nsec = 6
		}
		this.b = append(this.b, NanoSecSeparator)
		this.IntegerAlignRight(date_time.Nanosecond(), nsec)
	}
	this.b = append(this.b, SpaceSeparator)
}

func (this *Log) Level(level Level) {
	this.b = append(this.b, byte(level))
	this.b = append(this.b, SpaceSeparator)
}

func (this *Log) FilePathLine(skip int, fileLine FileLine) {
	if fileLine != FileLineDisable {
		_, f, l, o := runtime.Caller(skip + 1)
		if !o {
			this.b = append(this.b, unknownFileLine...)
		} else {
			if FileLineName == fileLine {
				for i := len(f) - 1; i >= 0; i-- {
					if f[i] == '/' {
						this.b = append(this.b, f[i+1:]...)
						this.b = append(this.b, FileLineSeparator)
						this.Integer(l)
						break
					}
				}
			} else {
				this.b = append(this.b, f...)
				this.b = append(this.b, FileLineSeparator)
				this.Integer(l)
			}
		}
	}
	this.b = append(this.b, FileLineSeparator)
	this.b = append(this.b, SpaceSeparator)
}

func (this *Log) String(s string) {
	this.b = append(this.b, s...)
}

func (this *Log) Write(b []byte) (int, error) {
	this.b = append(this.b, b...)
	return len(b), nil
}

func (this *Log) PrintBytes(writer io.Writer, level Level, skip int, fileLine FileLine, log []byte) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	this.b = append(this.b, log...)
	this.EndLine()
	return writer.Write(this.b)
}

func (this *Log) Print(writer io.Writer, level Level, skip int, fileLine FileLine, log string) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	this.String(log)
	this.EndLine()
	return writer.Write(this.b)
}

func (this *Log) Printf(writer io.Writer, level Level, skip int, fileLine FileLine, format string, a ... interface{}) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	fmt.Fprintf(this, format, a...)
	this.EndLine()
	return writer.Write(this.b)
}

func (this *Log) Sprint(writer io.Writer, level Level, skip int, fileLine FileLine, a ... interface{}) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	fmt.Fprint(this, a...)
	this.EndLine()
	return writer.Write(this.b)
}

func (this *Log) D(writer io.Writer, fileLine FileLine, log string) (int, error) {
	return this.Print(writer, LevelDebug, 1, fileLine, log)
}

func (this *Log) I(writer io.Writer, fileLine FileLine, log string) (int, error) {
	return this.Print(writer, LevelInfo, 1, fileLine, log)
}

func (this *Log) W(writer io.Writer, fileLine FileLine, log string) (int, error) {
	return this.Print(writer, LevelWarn, 1, fileLine, log)
}

func (this *Log) E(writer io.Writer, fileLine FileLine, log string) (int, error) {
	return this.Print(writer, LevelError, 1, fileLine, log)
}

func (this *Log) Stack(full, line bool) {
	i1 := len(this.b)
	i2 := 0
	n := 0
	for {
		this.b = this.b[:cap(this.b)]
		n = runtime.Stack(this.b[i1:], full)
		i2 = i1 + n
		if i2 < len(this.b) {
			this.b = this.b[:i2]
			break
		}
		this.b = append(this.b, make([]byte, 1024)...)
	}
	panic_line := false
	// 简化一下，只保留文件路径
	n = i2 - 1
	m := i1
Loop:
	for i := i1; i < n; {
		// 找出每一行
		if this.b[i] == '\t' {
			this.b[i] = ' '
			// 路径开始
			i++
			i1 = i
			for ; i < n; i++ {
				if this.b[i] == ' ' {
					// 路径结束
					i2 = i + 1
					// 只要一行
					if line {
						if !panic_line {
							panic_line = bytes.Contains(this.b[i1:i2], panicLine)
						} else {
							m += copy(this.b[m:], this.b[i1:i2])
							break Loop
						}
					} else {
						if !full {
							// 只要panic之后的路径
							if !panic_line {
								panic_line = bytes.Contains(this.b[i1:i2], panicLine)
							} else {
								m += copy(this.b[m:], this.b[i1:i2])
							}
						} else {
							m += copy(this.b[m:], this.b[i1:i2])
						}
					}
					break
				}
			}
		}
		i++
	}
	this.b = this.b[:m]
}

func Get() *Log {
	return logPool.Get().(*Log)
}

func Put(l *Log) {
	logPool.Put(l)
}

func Print(writer io.Writer, level Level, skip int, fileLine FileLine, log string) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Print(writer, level, skip+1, fileLine, log)
	logPool.Put(l)
	return n, e
}

func Printf(writer io.Writer, level Level, skip int, fileLine FileLine, format string, a ... interface{}) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Printf(writer, level, skip+1, fileLine, format, a...)
	logPool.Put(l)
	return n, e
}

func Sprint(writer io.Writer, level Level, skip int, fileLine FileLine, a ... interface{}) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Sprint(writer, level, skip+1, fileLine, a...)
	logPool.Put(l)
	return n, e
}

// 只输出panic的堆栈的路径，其他堆栈信息不输出
// 函数内部会调用recover()
// 参数
// writer: 输出
// full: 输出所有的堆栈，或者只输出panic之后的
// line: 只输出发生panic的那一行，line为true，full无效
// cb: 输出完后回调
// 返回
// 是否发生的panic
func Recover(writer io.Writer, full, line bool, cb func()) bool {
	// recover
	o := RecoverValue(writer, full, line, recover())
	// 回调函数
	if cb != nil {
		cb()
	}
	return o
}

// 只输出panic的堆栈的路径，其他堆栈信息不输出
// 参数
// writer: 输出
// full: 输出所有的堆栈，或者只输出panic之后的
// line: 只输出发生panic的那一行，line为true，full无效
// re: 在函数外调用recover()的值
// 返回
// 是否发生的panic
func RecoverValue(writer io.Writer, full, line bool, re interface{}) bool {
	if re != nil {
		// 获取Log
		l := logPool.Get().(*Log)
		l.Reset()
		if e, o := re.(*panicError); o {
			// 时间
			l.DateTime(6)
			// 级别，recover
			l.Level(LevelRecover)
			// 如果是log.Error，可以提升性能
			l.b = append(l.b, e.File...)
			l.b = append(l.b, FileLineSeparator)
			l.Integer(e.Line)
			l.b = append(l.b, FileLineSeparator)
			l.b = append(l.b, SpaceSeparator)
			l.String(e.Log)
			errPool.Put(e)
		} else {
			// 时间
			l.DateTime(6)
			// 级别，panic
			l.Level(LevelPanic)
			// 信息
			l.String(fmt.Sprint(re))
			l.Byte(SpaceSeparator)
			// 不是log.Error，从堆栈找到panic的行
			l.Stack(full, line)
		}
		l.EndLine()
		writer.Write(l.b)
		logPool.Put(l)
		return true
	}
	return false
}

type panicError struct {
	File string // 文件路径
	Line int    // 文件行
	Log  string // 信息
}

func (this *panicError) Error() string {
	return this.Log
}

func Panic(log string) {
	panic(newError(1, log))
}

func CheckError(e error) {
	if e != nil {
		panic(newError(1, e.Error()))
	}
}

func newError(skip int, log string) error {
	// 获取Error
	e := errPool.Get().(*panicError)
	_, f, l, o := runtime.Caller(skip + 1)
	if o {
		e.File = f
		e.Line = l
	} else {
		e.File = "???"
		e.Line = -1
	}
	e.Log = log
	return e
}

var (
	defaultWriter io.Writer = os.Stdout
)

func SetDefaultWriter(w io.Writer) {
	if nil != w {
		defaultWriter = w
	}
}

func Debug(a ... interface{}) {
	Sprint(os.Stderr, LevelDebug, 1, FileLineFullPath, a...)
}

func Info(a ... interface{}) {
	Sprint(os.Stderr, LevelInfo, 1, FileLineFullPath, a...)
}

func Warning(a ... interface{}) {
	Sprint(os.Stderr, LevelWarn, 1, FileLineFullPath, a...)
}

func Error(a ... interface{}) {
	Sprint(os.Stderr, LevelError, 1, FileLineFullPath, a...)
}
