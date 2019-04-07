package log

/*
1.I wanna print personal date format like 2018-01-01 00:00:00
2.I wanna print personal panic stack file line
*/

import (
	"time"
	"io"
	"fmt"
	"runtime"
	"bytes"
	"os"
	"sync"
)

type Level int

const (
	LevelDebug  Level = iota
	LevelInfo
	LevelWarn
	LevelError
	_LevelPanic
)

var (
	unknownFile            = "???"
	unknownStackLine       = [][]byte{[]byte("???:-1")}
	newline                = []byte("\n")
	space                  = []byte(" ")
	startHeader            = []byte("[")
	endHeader              = []byte("] ")
	colon                  = []byte(":")
	DateSeparator     byte = '-'
	TimeSeparator     byte = ':'
	DateTimeSeparator byte = ' '
	NanoSecSeparator  byte = '.'
	stdLogger              = NewStdLogger(LevelDebug, true)
)

func DefaultStdLogger() Logger {
	return stdLogger
}

type panicInfo struct {
	f string
	l int
	a interface{}
}

func (this *panicInfo) String() string {
	return fmt.Sprint(this.a)
}

func formatAlignInteger(b []byte, i int) {
	n := len(b) - 1
	j := i / 10
	for i >= 10 {
		b[n] = byte('0' + i - j*10)
		n--
		i = j
		j = i / 10
	}
	if n >= 0 {
		b[n] = byte('0' + i)
		n--
		for n >= 0 {
			b[n] = byte('0')
			n--
		}
	}
}

func formatInteger(b []byte, i int) int {
	n := 0
	for i > 0 {
		b[n] = byte('0' + i%10)
		i /= 10
		n++
	}
	k := n - 1
	m := byte(0)
	for j := 0; j < n; j++ {
		m = b[j]
		b[j] = b[k]
		b[k] = m
		k--
		if j >= k {
			break
		}
	}
	return n
}

func getStack() [][]byte {
	lines := make([][]byte, 0)
	stack := make([]byte, 128)
	for {
		n := runtime.Stack(stack, false)
		if n < len(stack) {
			stack = stack[:n]
			break
		}
		stack = make([]byte, len(stack)+128)
	}
	for len(stack) > 0 {
		i := bytes.IndexByte(stack, '\n')
		if i < 0 {
			return unknownStackLine
		}
		line := stack[:i]
		stack = stack[i+1:]
		if bytes.Contains(line, []byte("/runtime/panic.go")) {
			for len(stack) > 0 {
				i = bytes.IndexByte(stack, '\n')
				if i < 0 {
					return unknownStackLine
				}
				line = stack[:i]
				stack = stack[i+1:]
				if line[0] == '\t' {
					j := 1
					for ; j < len(line); j++ {
						if line[j] == ' ' {
							break
						}
					}
					lines = append(lines, line[1:j])
				}
			}
			return lines
		}
	}
	return unknownStackLine
}

func ParseLevel(s string) Level {
	switch s {
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelDebug
	}
}

func Print(w io.Writer, l Level, s bool, d int, i string) {
	w.Write(startHeader)
	printTimeAndLevel(w, l)
	if s {
		_, f, l, o := runtime.Caller(d + 1)
		if !o {
			printFileLine(w, &unknownFile, 0)
		} else {
			printFileLine(w, &f, l)
		}
	}
	w.Write(endHeader)
	io.WriteString(w, i)
	w.Write(newline)
}

func Printf(w io.Writer, l Level, s bool, d int, f string, a ... interface{}) {
	Print(w, l, s, d+1, fmt.Sprintf(f, a...))
}

func printTimeAndLevel(w io.Writer, l Level) {
	t := time.Now()
	var buf [31]byte
	b := buf[:]
	year, month, day := t.Date()
	formatAlignInteger(b[0:4], year)
	b[4] = DateSeparator
	formatAlignInteger(b[5:7], int(month))
	b[7] = DateSeparator
	formatAlignInteger(b[8:10], day)
	b[10] = DateTimeSeparator
	hour, min, sec := t.Clock()
	formatAlignInteger(b[11:13], hour)
	b[13] = TimeSeparator
	formatAlignInteger(b[14:16], min)
	b[16] = TimeSeparator
	formatAlignInteger(b[17:19], sec)
	b[19] = NanoSecSeparator
	formatAlignInteger(b[20:29], t.Nanosecond())
	b[29] = ' '
	switch l {
	case LevelDebug:
		b[30] = 'd'
	case LevelInfo:
		b[30] = 'i'
	case LevelWarn:
		b[30] = 'w'
	case LevelError:
		b[30] = 'e'
	case _LevelPanic:
		b[30] = 'p'
	default:
		b[30] = 'n'
	}
	w.Write(b)
}

func printFileLine(w io.Writer, s *string, l int) {
	w.Write(space)
	io.WriteString(w, *s)
	var buf [22]byte
	buf[0] = ':'
	b := buf[:]
	n := formatInteger(b[1:], l)
	n++
	b[n] = ':'
	w.Write(b[:n+1])
}

func Recover(w io.Writer, r interface{}) bool {
	if nil == r {
		return false
	}
	w.Write(startHeader)
	printTimeAndLevel(w, _LevelPanic)
	switch r.(type) {
	case *panicInfo:
		info := r.(*panicInfo)
		printFileLine(w, &info.f, info.l)
		w.Write(endHeader)
		text := fmt.Sprint(info.a)
		io.WriteString(w, text)
	default:
		w.Write(space)
		stacks := getStack()
		if len(stacks) > 0 {
			w.Write(stacks[0])
			w.Write(colon)
		}
		for i := 1; i < len(stacks); i++ {
			w.Write(space)
			w.Write(stacks[i])
			w.Write(colon)
		}
		w.Write(endHeader)
		text := fmt.Sprint(r)
		io.WriteString(w, text)
	}
	w.Write(newline)
	return true
}

func Panic(a interface{}) {
	if nil != a {
		ee := new(panicInfo)
		ee.a = a
		_, ee.f, ee.l, _ = runtime.Caller(1)
		panic(ee)
	}
}

func Panicf(f string, a ... interface{}) {
	ee := new(panicInfo)
	ee.a = fmt.Sprintf(f, a...)
	_, ee.f, ee.l, _ = runtime.Caller(1)
	panic(ee)
}

type Logger interface {
	Print(l Level, d int, s string)
	Printf(l Level, d int, f string, a ...interface{})
	Recover(r interface{}) bool
	SetLevel(l Level)
	io.Closer
}

type StdLogger struct {
	mux   sync.Mutex
	stack bool
	level Level
}

func (this *StdLogger) Print(l Level, d int, s string) {
	if l >= this.level {
		this.mux.Lock()
		Print(os.Stderr, l, this.stack, d+1, s)
		this.mux.Unlock()
	}
}

func (this *StdLogger) Printf(l Level, d int, f string, a ...interface{}) {
	if l >= this.level {
		this.mux.Lock()
		Printf(os.Stderr, l, this.stack, d+1, f, a...)
		this.mux.Unlock()
	}
}

func (this *StdLogger) Recover(r interface{}) bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return Recover(os.Stderr, r)
}

func (this *StdLogger) Close() error {
	return nil
}

func (this *StdLogger) SetLevel(l Level) {
	this.level = l
}

func NewStdLogger(level Level, stack bool) *StdLogger {
	return &StdLogger{level: level, stack: stack}
}
