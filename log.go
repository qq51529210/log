package log

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

const (
	LEVEL_DEBUG  Level = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_ERROR
	_LEVEL_PANIC
)

type Level int

type panicInfo struct {
	file  string
	line  int
	value interface{}
	time  time.Time
}

var (
	date_dir_fmt  = "20060102"
	time_file_fmt = "150405.999999999"
	levelFmt      = [][]byte{
		[]byte(string(" [DEBUG] ")),
		[]byte(string(" [WARN] ")),
		[]byte(string(" [INFO] ")),
		[]byte(string(" [ERROR] ")),
		[]byte(string(" [RECOVER] ")),
		[]byte(string(" [PANIC] ")),
	}
	unknownFileLine = []byte("???:-1")
	endLine         = []byte("\n")
	_log            = &logger{
		std:           true,
		valid:         false,
		stack:         make([]byte, 4096),
		fmtDateTime:   newFmtDateTime(),
		fmtLineNo:     newFmtLineNo(),
		panicFileLine: newPanicFileLine(),
	}
)

type panicFileLine []byte

func newPanicFileLine() panicFileLine {
	return make([]byte, 32)
}

func (this panicFileLine) Find() []byte {
	var stack []byte
	for {
		n := runtime.Stack(this, false)
		if n < len(this) {
			stack = this[:n]
			break
		}
		this = make([]byte, len(this)+1024)
	}
	for len(stack) > 0 {
		i := bytes.IndexByte(stack, '\n')
		if i < 0 {
			return unknownFileLine
		}
		line := stack[:i]
		stack = stack[i+1:]
		if bytes.Contains(line, []byte("/runtime/panic.go")) {
			for j := 0; j < 2; j++ {
				i = bytes.IndexByte(stack, '\n')
				if i < 0 {
					return unknownFileLine
				}
				line = stack[:i]
				stack = stack[i+1:]
			}
			for j := 0; j < len(line); j++ {
				if line[j] != ' ' && line[j] != '\t' {
					//line = line[j:]
					//for k := len(line) - 1; k >= 0; k-- {
					//	if line[k] == '/' || line[k] == '\\' {
					//		return line[k+1:]
					//	}
					//}
					return line[j:]
				}
			}
			break
		}
	}
	return unknownFileLine
}

type fmtLineNo []byte

func newFmtLineNo() fmtLineNo {
	return make([]byte, 32)
}

func (this fmtLineNo) Fmt(line int) int {
	this[0] = ':'
	n := fmtInt2(this[1:], line) + 1
	return n
}

type fmtDateTime []byte

func newFmtDateTime() fmtDateTime {
	return make([]byte, 30)
}

func (this fmtDateTime) Fmt(now time.Time) {
	year, month, day := now.Date()
	fmtInt(this[0:4], year)
	this[4] = '-'
	fmtInt(this[5:7], int(month))
	this[7] = '-'
	fmtInt(this[8:10], day)
	this[10] = ' '
	hour, min, sec := now.Clock()
	fmtInt(this[11:13], hour)
	this[13] = ':'
	fmtInt(this[14:16], min)
	this[16] = ':'
	fmtInt(this[17:19], sec)
	this[19] = '.'
	fmtInt(this[20:29], now.Nanosecond())
	this[29] = ' '
}

func fmtInt(b []byte, i int) {
	n := len(b) - 1
	for i >= 10 {
		j := i / 10
		b[n] = byte('0' + i - j*10)
		n--
		i = j
	}
	b[n] = byte('0' + i)
	n--
	for n >= 0 {
		b[n] = byte('0')
		n--
	}
}

func fmtInt2(b []byte, i int) int {
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

func Open(dir string, size, day, recent int, std bool) Logger {
	l := &logger{
		dir:           dir,
		size:          size,
		day:           day,
		recent:        recent,
		std:           std,
		valid:         true,
		stack:         make([]byte, 4096),
		fmtDateTime:   newFmtDateTime(),
		fmtLineNo:     newFmtLineNo(),
		panicFileLine: newPanicFileLine(),
	}

	if l.size < 1 {
		l.size = 1024 * 1024
	}
	if l.day < 1 {
		l.day = 7
	}

	if l.dir != "" {
		l.Add(1)
		go l.syncLoop()
	}
	return l
}

func Printf(level Level, skip int, format string, value ...interface{}) {
	text := fmt.Sprintf(format, value...)
	_log.print0(level, skip+2, &text)
}

func Print(level Level, skip int, text string) {
	_log.print0(level, skip+2, &text)
}

func Panic(skip int, value interface{}) {
	if nil == value {
		return
	}
	o := false
	info := &panicInfo{time: time.Now(), value: value}
	_, info.file, info.line, o = runtime.Caller(skip + 1)
	if !o {
		info.file = "???"
		info.line = -1
	}
	panic(info)
}
