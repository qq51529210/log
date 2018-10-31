package log

/*
1.I wanna print personal date format like 2018-01-01 00:00:00
2.I wanna print personal panic stack file line
*/

import (
	"strings"
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

type StackLine int

const (
	StackLineNil  StackLine = iota
	StackLineFile
	StackLinePath
)

var (
	unknownFile            = []byte("???")
	unknownStackLine       = [][]byte{[]byte("???:-1")}
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

func Print(w io.Writer, l Level, f StackLine, d int, s string) {
	t := time.Now()
	var buf [32]byte
	b := buf[:]
	printTimeAndLevel(b, &t, l)
	w.Write(b)
	switch f {
	case StackLineFile:
		_, f, l, o := runtime.Caller(d + 1)
		if !o {
			w.Write(unknownFile)
			w.Write(b[:printInt(b, 0)])
		} else {
			n := strings.LastIndex(f, "/")
			if n < 0 {
				break
			}
			w.Write(unsafeBytesFromString(&f)[n:])
			w.Write(b[:printInt(b, l)])
		}
	case StackLinePath:
		_, f, l, o := runtime.Caller(d + 1)
		if !o {
			w.Write(unknownFile)
			w.Write(b[:printInt(b, 0)])
		} else {
			w.Write(unsafeBytesFromString(&f))
			w.Write(b[:printInt(b, l)])
		}
	default:
	}
	w.Write(unsafeBytesFromString(&s))
	w.Write(newline)
}

func printTimeAndLevel(b []byte, t *time.Time, l Level) {
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
	case _LevelStack:
		b[30] = 's'
	default:
		b[30] = 'n'
	}
	b[31] = ' '
}

func printInt(b []byte, l int) int {
	b[0] = ':'
	n := formatInteger(b[1:], l)
	n++
	b[n] = ' '
	n++
	return n
}

func Recover(w io.Writer, r interface{}) bool {
	if nil == r {
		return false
	}
	var buf [32]byte
	b := buf[:]
	switch r.(type) {
	case *panicInfo:
		info := r.(*panicInfo)
		printTimeAndLevel(b, &info.t, _LevelPanic)
		w.Write(unsafeBytesFromString(&info.f))
		w.Write(b[:printInt(b, info.l)])
		text := fmt.Sprint(info.e.Error())
		w.Write(unsafeBytesFromString(&text))
	default:
		now := time.Now()
		printTimeAndLevel(b, &now, _LevelStack)
		stacks := getStack()
		if len(stacks) > 0 {
			w.Write(stacks[0])
			w.Write(space)
		}
		text := fmt.Sprint(r)
		w.Write(unsafeBytesFromString(&text))
		w.Write(space)
		for i := 1; i < len(stacks); i++ {
			w.Write(stacks[i])
			w.Write(space)
		}
	}
	w.Write(newline)
	return true
}
