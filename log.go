package log

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Level int

const (
	LEVEL_DEBUG   = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_ERROR
	LEVEL_RECOVER
)

var (
	date_dir_fmt  = "20060102"
	time_file_fmt = "150405.999999999"
	level_fmt     = [][]byte{
		[]byte(string("[d")),
	}
)

type Error struct {
}

type Log interface {
	Print(Level, int, string)
	//Panic(e error)
	//Recover(interface{})
	//io.Closer
}

//func Open(dir string, level string, size, day int, dur time.Duration) Log {
func Open() Log {
	return &log{buf: make([]byte, 30)}
}

type log struct {
	sync.Mutex
	buf  []byte
	data bytes.Buffer
}

func (this *log) Print(level Level, depth int, s string) {
	this.Lock()
	// date time
	now := time.Now()
	year, month, day := now.Date()
	fmtInt(this.buf[0:4], year)
	this.buf[4] = '-'
	fmtInt(this.buf[5:7], int(month))
	this.buf[7] = '-'
	fmtInt(this.buf[8:10], day)
	this.buf[10] = ' '
	hour, min, sec := now.Clock()
	fmtInt(this.buf[11:13], hour)
	this.buf[13] = ':'
	fmtInt(this.buf[14:16], min)
	this.buf[16] = ':'
	fmtInt(this.buf[17:19], sec)
	this.buf[19] = '.'
	fmtInt(this.buf[20:29], now.Nanosecond())
	this.buf[29] = ' '
	this.data.Write(this.buf[:30])
	// file line
	_, f, l, o := runtime.Caller(depth)
	if !o {
		f = "???"
		l = -1
	}
	this.data.WriteString(f)
	this.buf[0] = ':'
	n := fmtInt2(this.buf[1:], l) + 1
	this.buf[n] = ':'
	n++
	this.buf[n] = ' '
	n++
	this.buf[n] = '['
	n++
	switch level {
	case LEVEL_DEBUG:
		this.buf[n] = 'd'
	case LEVEL_WARN:
		this.buf[n] = 'w'
	case LEVEL_INFO:
		this.buf[n] = 'i'
	case LEVEL_ERROR:
		this.buf[n] = 'e'
	case LEVEL_RECOVER:
		this.buf[n] = 'r'
	}
	n++
	this.buf[n] = ']'
	n++
	this.buf[n] = ' '
	n++
	this.data.Write(this.buf[:n])
	this.data.WriteString(s)
	this.Unlock()

	fmt.Println(string(this.data.Bytes()))
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
	for i >= 10 {
		j := i
		for j >= 10 {
			j /= 10
		}
		b[n] = byte('0' + j)
		n++
		i %= 10 * n
	}
	b[n] = byte('0' + i)
	n++
	return n
}
