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
	LevelDebug       Level = 'd'
	LevelInfo        Level = 'i'
	LevelWarn        Level = 'w'
	LevelError       Level = 'e'
	LevelPanic       Level = 'p'
	MaxIntegerLength       = 20
)

type FileLine byte

const (
	FileLineDisable  FileLine = iota
	FileLineFullPath
	FileLineName
)

var (
	logPool                = sync.Pool{}
	SpaceSeparator    byte = ' '
	DateSeparator     byte = '-'
	TimeSeparator     byte = ':'
	NanoSecSeparator  byte = '.'
	FileLineSeparator byte = ':'
	unknownFileLine        = []byte("???:-1:")
)

func init() {
	logPool.New = func() interface{} {
		return &Log{}
	}
}

type Log struct {
	b []byte
	n int
}

func (this *Log) Reset() {
	this.n = 0
	this.Grow(34)
}

func (this *Log) Grow(n int) {
	left := len(this.b) - this.n
	if left < n {
		this.b = append(this.b, make([]byte, n-left)...)
	}
}

func (this *Log) integerAlignRight(value, length int) {
	max := this.n + length
	// 值是倒转的
	m := this.n
	for value > 0 {
		this.b[m] = byte('0' + value%10)
		value /= 10
		m++
		length--
	}
	// 先反转
	c := byte(0)
	i := this.n
	this.n = m
	m--
	for i < m {
		c = this.b[i]
		this.b[i] = this.b[m]
		this.b[m] = c
		m--
		i++
	}
	// 再后面补0
	for length > 0 {
		this.b[m] = '0'
		m++
		length--
	}
	if this.n > max {
		this.n = max
	}
}

func (this *Log) integerAlignLeft(value, length int) {
	max := this.n + length
	// 值是倒转的
	m := this.n
	for value > 0 {
		this.b[m] = byte('0' + value%10)
		value /= 10
		m++
		length--
	}
	// 继续在后面补0
	for length > 0 {
		this.b[m] = '0'
		m++
		length--
	}
	// 反转
	c := byte(0)
	i := this.n
	this.n = m
	m--
	for i < m {
		c = this.b[i]
		this.b[i] = this.b[m]
		this.b[m] = c
		m--
		i++
	}
	if this.n > max {
		this.n = max
	}
}

func (this *Log) integer(n int) {
	m := this.n
	for n > 0 {
		this.b[m] = byte('0' + n%10)
		n /= 10
		m++
	}
	c := byte(0)
	i := this.n
	this.n = m
	m--
	for i < m {
		c = this.b[i]
		this.b[i] = this.b[m]
		this.b[m] = c
		m--
		i++
	}
}

func (this *Log) DateTime(nsec int) {
	// 2019-04-25 00:25:43:25.256131000_
	date_time := time.Now()
	year, month, day := date_time.Date()
	hour, minute, second := date_time.Clock()
	// fmt
	this.integerAlignLeft(year, 4)
	this.b[this.n] = DateSeparator
	this.n++
	this.integerAlignLeft(int(month), 2)
	this.b[this.n] = DateSeparator
	this.n++
	this.integerAlignLeft(day, 2)
	this.b[this.n] = SpaceSeparator
	this.n++
	this.integerAlignLeft(hour, 2)
	this.b[this.n] = TimeSeparator
	this.n++
	this.integerAlignLeft(minute, 2)
	this.b[this.n] = TimeSeparator
	this.n++
	this.integerAlignLeft(second, 2)
	if nsec > 0 {
		if nsec > 9 {
			nsec = 9
		}
		this.b[this.n] = NanoSecSeparator
		this.n++
		this.integerAlignRight(date_time.Nanosecond(), nsec)
	}
	this.b[this.n] = SpaceSeparator
	this.n++
}

func (this *Log) Byte(c byte) {
	this.b[this.n] = c
	this.n++
}

func (this *Log) Level(level Level) {
	this.b[this.n] = byte(level)
	this.n++
	this.b[this.n] = SpaceSeparator
	this.n++
}

func (this *Log) FilePathLine(skip int, fileLine FileLine) {
	if fileLine == FileLineDisable {
		return
	}
	_, f, l, o := runtime.Caller(skip + 1)
	if !o {
		this.unknownFilePathLine()
	} else {
		if FileLineFullPath == fileLine {
			for i := len(f) - 1; i >= 0; i-- {
				if f[i] == '/' {
					this.filePathLine(f[i+1:], l)
					break
				}
			}
		} else {
			this.filePathLine(f, l)
		}
	}
	this.b[this.n] = SpaceSeparator
	this.n++
}

func (this *Log) unknownFilePathLine() {
	this.Grow(len(unknownFileLine) + 1)
	this.n += copy(this.b[:this.n], unknownFileLine)
}

func (this *Log) filePathLine(f string, l int) {
	this.Grow(len(f) + MaxIntegerLength + 1)
	this.n += copy(this.b[this.n:], f)
	this.Byte(FileLineSeparator)
	this.integer(l)
}

func (this *Log) String(s string) {
	this.Grow(len(s))
	this.n += copy(this.b[this.n:], s)
}

func (this *Log) Write(b []byte) (int, error) {
	this.Grow(len(b))
	this.n += copy(this.b[this.n:], b)
	return len(b), nil
}

func (this *Log) Print(writer io.Writer, level Level, skip int, fileLine FileLine, log string) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	this.String(log)
	this.Grow(1)
	this.Byte('\n')
	return writer.Write(this.b[:this.n])
}

func (this *Log) Printf(writer io.Writer, level Level, skip int, fileLine FileLine, format string, a ... interface{}) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(level)
	this.FilePathLine(skip+1, fileLine)
	fmt.Fprintf(this, format, a...)
	this.Grow(1)
	this.Byte('\n')
	return writer.Write(this.b[:this.n])
}

func Print(writer io.Writer, level Level, skip int, fileLine FileLine, log string) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Print(writer, level, skip, fileLine, log)
	logPool.Put(l)
	return n, e

	//l := logPool.Get().(*Log)
	//l.Reset()
	//l.DateTime(6)
	//l.Level(level)
	//l.FilePathLine(skip+1, fileLine)
	//l.String(log)
	//l.grow(1)
	//l.Byte('\n')
	//n, e := writer.Write(l.b[:l.n])
	//logPool.Put(l)
	//return n, e
}

func Printf(writer io.Writer, level Level, skip int, fileLine FileLine, format string, a ... interface{}) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.Printf(writer, level, skip, fileLine, format, a...)
	logPool.Put(l)
	return n, e

	//l := logPool.Get().(*Log)
	//l.Reset()
	//l.DateTime(6)
	//l.Level(level)
	//l.FilePathLine(skip+1, fileLine)
	//fmt.Fprintf(f, format, a...)
	//l.grow(1)
	//l.Byte('\n')
	//n, e := writer.Write(l.b[:l.n])
	//logPool.Put(l)
	//return n, e
}

func Recover(writer io.Writer) {
	re := recover()
	if re == nil {
		return
	}

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
//	formatAligninteger(b[0:4], year)
//	b[4] = DateSeparator
//	formatAligninteger(b[5:7], int(month))
//	b[7] = DateSeparator
//	formatAligninteger(b[8:10], day)
//	b[10] = DateTimeSeparator
//	hour, min, sec := t.Clock()
//	formatAligninteger(b[11:13], hour)
//	b[13] = TimeSeparator
//	formatAligninteger(b[14:16], min)
//	b[16] = TimeSeparator
//	formatAligninteger(b[17:19], sec)
//	b[19] = NanoSecSeparator
//	formatAligninteger(b[20:29], t.Nanosecond())
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
//	n := formatinteger(b[1:], l)
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
