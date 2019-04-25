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
	re := recover()
	if re == nil {
		return
	}
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
	// 回调函数
	if cb != nil {
		cb()
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
	// 获取Error
	e := errPool.Get().(*Error)
	_, f, l, o := runtime.Caller(1)
	if o {
		e.File = f
		e.Line = l
	} else {
		e.File = "???"
		e.Line = -1
	}
	e.Log = log
	// panic
	panic(e)
}

//
///*
//1.I wanna print personal date format like 2018-01-01 00:00:00
//2.I wanna print personal panic stack file line
//*/
//
//import (
//	"time"
//	"io"
//	"fmt"
//	"runtime"
//	"bytes"
//	"os"
//	"sync"
//)
//
//type Level int
//
//var (
//	unknownFile            = "???"
//	unknownStackLine       = [][]byte{[]byte("???:-1")}
//	newline                = []byte("\n")
//	space                  = []byte(" ")
//	startHeader            = []byte("[")
//	endHeader              = []byte("] ")
//	colon                  = []byte(":")
//	DateSeparator     byte = '-'
//	TimeSeparator     byte = ':'
//	DateTimeSeparator byte = ' '
//	NanoSecSeparator  byte = '.'
//	stdLogger              = NewStdLogger(LevelDebug, true)
//)
//
//func DefaultStdLogger() Log {
//	return stdLogger
//}
//
//type panicInfo struct {
//	f string
//	l int
//	a interface{}
//}
//
//func (this *panicInfo) String() string {
//	return fmt.Sprint(this.a)
//}
//
//func getStack() [][]byte {
//	lines := make([][]byte, 0)
//	stack := make([]byte, 128)
//	for {
//		n := runtime.Stack(stack, false)
//		if n < len(stack) {
//			stack = stack[:n]
//			break
//		}
//		stack = make([]byte, len(stack)+128)
//	}
//	for len(stack) > 0 {
//		i := bytes.IndexByte(stack, '\n')
//		if i < 0 {
//			return unknownStackLine
//		}
//		line := stack[:i]
//		stack = stack[i+1:]
//		if bytes.Contains(line, []byte("/runtime/panic.go")) {
//			for len(stack) > 0 {
//				i = bytes.IndexByte(stack, '\n')
//				if i < 0 {
//					return unknownStackLine
//				}
//				line = stack[:i]
//				stack = stack[i+1:]
//				if line[0] == '\t' {
//					j := 1
//					for ; j < len(line); j++ {
//						if line[j] == ' ' {
//							break
//						}
//					}
//					lines = append(lines, line[1:j])
//				}
//			}
//			return lines
//		}
//	}
//	return unknownStackLine
//}
//
//func ParseLevel(s string) Level {
//	switch s {
//	case "info":
//		return LevelInfo
//	case "warn":
//		return LevelWarn
//	case "error":
//		return LevelError
//	default:
//		return LevelDebug
//	}
//}
//
//func Print(w io.Writer, l Level, s bool, d int, i string) {
//	w.Write(startHeader)
//	printTimeAndLevel(w, l)
//	if s {
//		_, f, l, o := runtime.Caller(d + 1)
//		if !o {
//			printFileLine(w, &unknownFile, 0)
//		} else {
//			printFileLine(w, &f, l)
//		}
//	}
//	w.Write(endHeader)
//	io.WriteString(w, i)
//	w.Write(newline)
//}
//
//func Printf(w io.Writer, l Level, s bool, d int, f string, a ... interface{}) {
//	Print(w, l, s, d+1, fmt.Sprintf(f, a...))
//}
//
//func printTimeAndLevel(w io.Writer, l Level) {
//	t := time.Now()
//	var buf [31]byte
//	b := buf[:]
//	year, month, day := t.Date()
//	formatAlignInteger(b[0:4], year)
//	b[4] = DateSeparator
//	formatAlignInteger(b[5:7], int(month))
//	b[7] = DateSeparator
//	formatAlignInteger(b[8:10], day)
//	b[10] = DateTimeSeparator
//	hour, min, sec := t.Clock()
//	formatAlignInteger(b[11:13], hour)
//	b[13] = TimeSeparator
//	formatAlignInteger(b[14:16], min)
//	b[16] = TimeSeparator
//	formatAlignInteger(b[17:19], sec)
//	b[19] = NanoSecSeparator
//	formatAlignInteger(b[20:29], t.Nanosecond())
//	b[29] = ' '
//	switch l {
//	case LevelDebug:
//		b[30] = 'd'
//	case LevelInfo:
//		b[30] = 'i'
//	case LevelWarn:
//		b[30] = 'w'
//	case LevelError:
//		b[30] = 'e'
//	case _LevelPanic:
//		b[30] = 'p'
//	default:
//		b[30] = 'n'
//	}
//	w.Write(b)
//}
//
//func printFileLine(w io.Writer, s *string, l int) {
//	w.Write(space)
//	io.WriteString(w, *s)
//	var buf [22]byte
//	buf[0] = ':'
//	b := buf[:]
//	n := formatInteger(b[1:], l)
//	n++
//	b[n] = ':'
//	w.Write(b[:n+1])
//}
//
//func Recover(w io.Writer, r interface{}) bool {
//	if nil == r {
//		return false
//	}
//	w.Write(startHeader)
//	printTimeAndLevel(w, _LevelPanic)
//	switch r.(type) {
//	case *panicInfo:
//		info := r.(*panicInfo)
//		printFileLine(w, &info.f, info.l)
//		w.Write(endHeader)
//		text := fmt.Sprint(info.a)
//		io.WriteString(w, text)
//	default:
//		w.Write(space)
//		stacks := getStack()
//		if len(stacks) > 0 {
//			w.Write(stacks[0])
//			w.Write(colon)
//		}
//		for i := 1; i < len(stacks); i++ {
//			w.Write(space)
//			w.Write(stacks[i])
//			w.Write(colon)
//		}
//		w.Write(endHeader)
//		text := fmt.Sprint(r)
//		io.WriteString(w, text)
//	}
//	w.Write(newline)
//	return true
//}
//
//func Panic(a interface{}) {
//	if nil != a {
//		ee := new(panicInfo)
//		ee.a = a
//		_, ee.f, ee.l, _ = runtime.Caller(1)
//		panic(ee)
//	}
//}
//
//func Panicf(f string, a ... interface{}) {
//	ee := new(panicInfo)
//	ee.a = fmt.Sprintf(f, a...)
//	_, ee.f, ee.l, _ = runtime.Caller(1)
//	panic(ee)
//}
//
//type Log interface {
//	Print(l Level, d int, s string)
//	Printf(l Level, d int, f string, a ...interface{})
//	Recover(r interface{}) bool
//	SetLevel(l Level)
//	io.Closer
//}
//
//type StdLogger struct {
//	mux   sync.Mutex
//	stack bool
//	level Level
//}
//
//func (this *StdLogger) Print(l Level, d int, s string) {
//	if l >= this.level {
//		this.mux.Lock()
//		Print(os.Stderr, l, this.stack, d+1, s)
//		this.mux.Unlock()
//	}
//}
//
//func (this *StdLogger) Printf(l Level, d int, f string, a ...interface{}) {
//	if l >= this.level {
//		this.mux.Lock()
//		Printf(os.Stderr, l, this.stack, d+1, f, a...)
//		this.mux.Unlock()
//	}
//}
//
//func (this *StdLogger) Recover(r interface{}) bool {
//	this.mux.Lock()
//	defer this.mux.Unlock()
//	return Recover(os.Stderr, r)
//}
//
//func (this *StdLogger) Close() error {
//	return nil
//}
//
//func (this *StdLogger) SetLevel(l Level) {
//	this.level = l
//}
//
//func NewStdLogger(level Level, stack bool) *StdLogger {
//	return &StdLogger{level: level, stack: stack}
//}
//
//// format
//// 2019-04-24 18:15:56.602572000
//// 29字节
//type TimeFormat [30]byte
