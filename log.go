package log

import (
	"fmt"
	"io"
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
		return &Error{}
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
	//this.String(fmt.Sprintf(format,a...))
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

func (this *Log) Stack() {
	i1 := len(this.b)
	i2 := 0
	n := 0
	for {
		this.b = append(this.b, make([]byte, 64)...)
		n = runtime.Stack(this.b[i1:], false)
		i2 = i1 + n
		if i2 < len(this.b) {
			this.b = this.b[:i2]
			break
		}
	}
	// 简化一下
	n = i2 - 1
	m := i1
	for i := i1; i < n; {
		// 找出每一行
		if this.b[i] == '\t' {
			this.b[i] = ' '
			// 开始
			i1 = i + 1
			i++
			for ; i < n; i++ {
				if this.b[i] == ' ' {
					i2 = i + 1
					m += copy(this.b[m:], this.b[i1:i2])
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

func Recover(writer io.Writer, cb func()) {
	// recover
	RecoverValue(writer, recover())
	// 回调函数
	if cb != nil {
		cb()
	}
}

func RecoverValue(writer io.Writer, re interface{}) {
	if re != nil {
		// 获取Log
		l := logPool.Get().(*Log)
		l.Reset()
		if e, o := re.(*Error); o {
			// 级别，recover
			l.Level(LevelRecover)
			// 时间
			l.DateTime(6)
			// 如果是log.Error，可以提升性能
			l.b = append(l.b, e.File...)
			l.b = append(l.b, FileLineSeparator)
			l.Integer(e.Line)
			l.b = append(l.b, FileLineSeparator)
			l.b = append(l.b, SpaceSeparator)
			l.String(e.Log)
			errPool.Put(e)
		} else {
			// 级别，panic
			l.Level(LevelPanic)
			// 时间
			l.DateTime(6)
			// 信息
			l.String(fmt.Sprint(re))
			l.Byte(SpaceSeparator)
			// 不是log.Error，从堆栈找到panic的行
			l.Stack()
		}
		l.EndLine()
		writer.Write(l.b)
		logPool.Put(l)
	}
}

type Error struct {
	File string // 文件路径
	Line int    // 文件行
	Log  string // 信息
}

func (this *Error) Error() string {
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
	e := errPool.Get().(*Error)
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
