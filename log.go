package log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	LEVEL_PANIC
)

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
)

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

func Open(dir string, size, day int, std bool) Log {
	l := &log{
		buf:   make([]byte, 30),
		dir:   dir,
		size:  size,
		day:   day,
		std:   std,
		valid: true,
		stack: make([]byte, 4096),
	}
	if l.dir != "" {
		l.Add(1)
		go l.syncLoop()
	}
	return l
}

func Panic(e error) {
	if nil == e {
		return
	}

	o := false
	err := &logError{time: time.Now(), err: e}
	_, err.file, err.line, o = runtime.Caller(1)
	if !o {
		err.file = "???"
		err.line = -1
	}

	panic(err)
}

type logError struct {
	file string
	line int
	err  error
	time time.Time
}

type Log interface {
	Print(Level, string)
	Recover()
	RecoverError(interface{})
	Close()
}

type log struct {
	sync.Mutex
	sync.WaitGroup
	dir   string
	size  int
	day   int
	std   bool
	buf   []byte
	data  bytes.Buffer
	valid bool
	stack []byte
	line  bytes.Buffer
	file  *os.File
}

func (this *log) print1(level Level, file, text *string, line int) {
	this.Lock()
	if this.std {
		// date time
		this.printDateTime(&this.line, time.Now())
		// file line
		this.printFileLine(&this.line, file, line)
		// level
		this.line.Write(levelFmt[level])
		// text
		this.line.WriteString(*text)
		// new line
		this.line.WriteByte('\n')
		// std
		os.Stderr.Write(this.line.Bytes())
		if "" != this.dir {
			this.data.Write(this.line.Bytes())
		}
		this.line.Reset()
	} else {
		if "" != this.dir {
			// date time
			this.printDateTime(&this.data, time.Now())
			// file line
			this.printFileLine(&this.data, file, line)
			// level
			this.data.Write(levelFmt[level])
			// text
			this.data.WriteString(*text)
			// new line
			this.data.WriteByte('\n')
		}
	}
	this.Unlock()
}

func (this *log) printDateTime(buf *bytes.Buffer, now time.Time) {
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
	buf.Write(this.buf[:30])
}

func (this *log) printFileLine(buf *bytes.Buffer, file *string, line int) {
	buf.WriteString(*file)
	this.buf[0] = ':'
	n := fmtInt2(this.buf[1:], line) + 1
	this.buf[n] = ':'
	n++
	buf.Write(this.buf[:n])
}

func (this *log) Print(level Level, text string) {
	_, f, l, o := runtime.Caller(1)
	if !o {
		f = "???"
		l = -1
	}
	this.print1(level, &f, &text, l)
}

func (this *log) findPanicFileLine() []byte {
	var stack []byte
	for {
		n := runtime.Stack(this.stack, false)
		if n < len(this.stack) {
			stack = this.stack[:n]
			break
		}
		this.stack = make([]byte, len(this.stack)+1024)
	}
	if len(stack) > 0 {
		for len(stack) > 0 {
			i := bytes.IndexByte(stack, '\n')
			if i < 0 {
				stack = unknownFileLine
				break
			}
			if bytes.Contains(stack[:i], []byte("/runtime/panic.go")) {
				stack = stack[i+1:]
				for i := 0; i < 2; i++ {
					i = bytes.IndexByte(stack, '\n')
					if i < 0 {
						stack = unknownFileLine
						break
					}
					stack = stack[i+1:]
				}
				stack = bytes.TrimLeftFunc(stack, func(r rune) bool {
					if r == '\t' || r == ' ' {
						return true
					}
					return false
				})
				i = bytes.IndexByte(stack, ' ')
				if i < 0 {
					stack = unknownFileLine
					break
				}
				stack = stack[:i]
				break
			}
			stack = stack[i+1:]
		}
	} else {
		stack = unknownFileLine
	}
	return stack
}

func (this *log) Recover() {
	this.RecoverError(recover())
}

func (this *log) RecoverError(re interface{}) {
	if nil == re {
		return
	}
	switch re.(type) {
	case *logError:
		err := re.(*logError)
		text := err.err.Error()
		this.print1(LEVEL_RECOVER, &err.file, &text, err.line)
	default:
		this.Lock()
		if this.std {
			// date time
			this.printDateTime(&this.line, time.Now())
			// file line
			this.line.Write(this.findPanicFileLine())
			// level
			this.line.Write(levelFmt[LEVEL_PANIC])
			// text
			this.line.WriteString(fmt.Sprintln(re))
			// std
			os.Stderr.Write(this.line.Bytes())
			if "" != this.dir {
				this.data.Write(this.line.Bytes())
			}
			this.line.Reset()
		} else {
			if "" != this.dir {
				// date time
				this.printDateTime(&this.data, time.Now())
				// file line
				this.data.Write(this.findPanicFileLine())
				// level
				this.data.Write(levelFmt[LEVEL_PANIC])
				// text
				this.data.WriteString(fmt.Sprintln(re))
			}
		}
		this.Unlock()
	}

}

func (this *log) Close() {
	this.Lock()
	if !this.valid {
		this.Unlock()
		return
	}
	this.valid = false
	this.Unlock()

	this.Wait()

	if nil != this.file {
		this.file.Close()
		this.file = nil
	}
}

func (this *log) syncLoop() {
	timer := time.NewTimer(time.Second)
	defer func() {
		recover()
		timer.Stop()
		this.Done()
	}()
	for this.valid {
		<-timer.C
		this.Lock()
		if nil == this.file {
			this.newFile()
		}
		if this.data.Len() > 0 {
			_, e := this.file.Write(this.data.Bytes())
			if nil == e {
				this.data.Reset()
				f, e := this.file.Stat()
				if nil == e {
					if int(f.Size()) >= this.size {
						this.file.Close()
						this.newFile()
					}
				} else {
					this.file.Close()
					this.newFile()
				}
			} else {
				this.newFile()
				_, e = this.file.Write(this.data.Bytes())
				if nil != e {
					this.file.Close()
					this.file = nil
					if this.data.Len() >= this.size {
						this.data.Reset()
					}
				} else {
					this.data.Reset()
				}
			}
		}
		this.Unlock()
		timer.Reset(time.Second)
	}
}

func (this *log) newFile() {
	now := time.Now()
	dir := filepath.Join(this.dir, now.Format(date_dir_fmt))
	e := os.MkdirAll(dir, os.ModePerm)
	if nil != e {
		os.Stderr.WriteString(e.Error())
		return
	}
	path := filepath.Join(dir, now.Format(time_file_fmt)+".log")
	f, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if nil != e {
		os.Stderr.WriteString(e.Error())
		return
	}
	this.file = f

	fs, e := ioutil.ReadDir(this.dir)
	if nil != e {
		os.Stderr.WriteString(e.Error())
		return
	}
	count := len(fs)
	if count <= this.day {
		return
	}
	for count > this.day {
		mt := time.Now()
		mt = time.Date(mt.Year(), mt.Month(), mt.Day(), 0, 0, 0, 0, time.Local)
		mi := -1
		for i := 0; i < len(fs); i++ {
			t, e := time.Parse(date_dir_fmt, fs[i].Name())
			if nil != e {
				count--
				e = os.RemoveAll(filepath.Join(this.dir, fs[i].Name()))
				if nil != e {
					os.Stderr.WriteString(e.Error())
				}
				continue
			}
			if t.Sub(mt) < 0 {
				mt = t
				mi = i
			}
		}
		if mi >= 0 {
			count--
			e = os.RemoveAll(filepath.Join(this.dir, fs[mi].Name()))
			if nil != e {
				os.Stderr.WriteString(e.Error())
			}
		}
	}
}
