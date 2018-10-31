package log

/*
1.I wanna print personal date format like 2018-01-01 00:00:00
2.I wanna print personal panic stack file line
*/

import (
	"time"
	"unsafe"
	"reflect"
	"io"
	"fmt"
	"runtime"
	"bytes"
)

type Level int

const (
	LevelDebug  Level = iota
	LevelInfo
	LevelWarn
	LevelError
	_LevelPanic
	_LevelStack
)

type FileLine int

const (
	FileLineNil  FileLine = iota
	FileLineFile
	FileLinePath
)

var (
	unknownFileLine        = [][]byte{[]byte("???:-1")}
	newline                = []byte("\n")
	space                  = []byte(" ")
	DateSeparator     byte = '-'
	TimeSeparator     byte = ':'
	DateTimeSeparator byte = ' '
	NanoSecSeparator  byte = '.'
)

type panicInfo struct {
	f string
	l int
	e error
	t time.Time
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

func unsafeBytesFromString(s *string) []byte {
	ss := (*reflect.StringHeader)(unsafe.Pointer(s))
	bb := reflect.SliceHeader{
		Data: ss.Data,
		Len:  ss.Len,
		Cap:  ss.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bb))
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
			return unknownFileLine
		}
		line := stack[:i]
		stack = stack[i+1:]
		if bytes.Contains(line, []byte("/runtime/panic.go")) {
			for len(stack) > 0 {
				i = bytes.IndexByte(stack, '\n')
				if i < 0 {
					return unknownFileLine
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
	return unknownFileLine
}

func Print(w io.Writer, l Level, f FileLine, d int, s string) {
	t := time.Now()
	PrintTime(w, &t)
	PrintLevel(w, l)
	switch f {
	case FileLineFile:
		_, f, l, o := runtime.Caller(d + 1)
		if !o {
			PrintFileLine(w, "???", -1)
		} else {
			PrintFileLine(w, f, l)
		}
	case FileLinePath:
		_, f, l, o := runtime.Caller(d + 1)
		if !o {
			PrintFileLine(w, "???", -1)
		} else {
			PrintFileLine(w, f, l)
		}
	default:
	}
	PrintText(w, s)
}

func PrintTime(w io.Writer, t *time.Time) {
	var buf [30]byte
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
	w.Write(b)
}

func PrintLevel(w io.Writer, l Level) {
	var buf [2]byte
	b := buf[:]
	switch l {
	case LevelDebug:
		b[0] = 'd'
	case LevelInfo:
		b[0] = 'i'
	case LevelWarn:
		b[0] = 'w'
	case LevelError:
		b[0] = 'e'
	case _LevelPanic:
		b[0] = 'p'
	case _LevelStack:
		b[0] = 's'
	default:
		b[0] = 'n'
	}
	b[1] = ' '
	w.Write(b)
}

func PrintFileLine(w io.Writer, f string, l int) {
	w.Write(unsafeBytesFromString(&f))
	var temp [22]byte
	p := temp[:]
	p[0] = ':'
	n := formatInteger(p[1:], l)
	n++
	p[n] = ' '
	n++
	w.Write(p[:n])
}

func PrintText(w io.Writer, s string) {
	w.Write(unsafeBytesFromString(&s))
	w.Write(newline)
}

func Recover(w io.Writer, r interface{}) bool {
	if nil == r {
		return false
	}
	switch r.(type) {
	case *panicInfo:
		info := r.(*panicInfo)
		text := fmt.Sprint(info.e.Error())
		PrintTime(w, &info.t)
		PrintLevel(w, _LevelPanic)
		PrintFileLine(w, info.f, info.l)
		PrintText(w, text)
	default:
		now := time.Now()
		text := fmt.Sprint(r)
		stacks := getStack()
		PrintTime(w, &now)
		PrintLevel(w, _LevelStack)
		if len(stacks) > 0 {
			w.Write(stacks[0])
			w.Write(space)
		}
		w.Write(unsafeBytesFromString(&text))
		w.Write(space)
		for i := 1; i < len(stacks); i++ {
			w.Write(stacks[i])
			w.Write(space)
		}
		w.Write(newline)
	}
	return true
}
