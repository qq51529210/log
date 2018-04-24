package log

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	LEVEL_DEBUG   Level = iota
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
	endLine         = []byte("\n")
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

func Open(dir string, size, day, recent int, std bool) Log {
	l := &log{
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

func Print(level Level, text string) {
	t := newFmtDateTime()
	t.Fmt(time.Now())
	os.Stderr.Write(t)
	_, f, l, o := runtime.Caller(1)
	if !o {
		f = "???"
		l = -1
	}
	os.Stderr.WriteString(f)
	no := newFmtLineNo()
	os.Stderr.Write(no[:no.Fmt(l)])
	os.Stderr.Write(levelFmt[level])
	os.Stderr.WriteString(text)
	os.Stderr.Write(endLine)
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

type fmtLineNo []byte

func newFmtLineNo() fmtLineNo {
	return make([]byte, 32)
}

func (this fmtLineNo) Fmt(line int) int {
	this[0] = ':'
	n := fmtInt2(this[1:], line) + 1
	this[n] = ':'
	n++
	return n
}

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

type Level int

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
	Recently() []string
}

type log struct {
	sync.RWMutex
	sync.WaitGroup
	dir    string
	size   int
	day    int
	recent int
	std    bool
	data   bytes.Buffer
	valid  bool
	stack  []byte
	line   bytes.Buffer
	file   *os.File
	list   list.List
	fmtDateTime
	fmtLineNo
	panicFileLine
}

func (this *log) Print(level Level, text string) {
	_, f, l, o := runtime.Caller(1)
	if !o {
		f = "???"
		l = -1
	}
	this.print1(level, &f, &text, l)
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
		text := fmt.Sprintln(re)
		file_line := this.panicFileLine.Find()
		this.Lock()
		if this.recent < 1 {
			if this.std {
				this.print3(&this.line, LEVEL_PANIC, file_line, &text)
				os.Stderr.Write(this.line.Bytes())
				if "" != this.dir {
					this.data.Write(this.line.Bytes())
				}
				this.line.Reset()
			} else {
				if "" != this.dir {
					this.print3(&this.data, LEVEL_PANIC, file_line, &text)
				}
			}
		} else {
			this.print3(&this.line, LEVEL_PANIC, file_line, &text)
			if this.list.Len() >= this.recent {
				this.list.Remove(this.list.Front())
			}
			this.list.PushBack(string(this.line.Bytes()))
			if this.std {
				os.Stderr.Write(this.line.Bytes())
			}
			if "" != this.dir {
				this.data.Write(this.line.Bytes())
			}
			this.line.Reset()
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

func (this *log) Recently() []string {
	recent := make([]string, 0)
	this.RLock()
	for ele := this.list.Front(); nil != ele; ele = ele.Next() {
		recent = append(recent, ele.Value.(string))
	}
	this.RUnlock()
	return recent
}

func (this *log) print1(level Level, file, text *string, line int) {
	this.Lock()
	if this.recent < 1 {
		if this.std {
			this.print2(&this.line, level, file, text, line)
			os.Stderr.Write(this.line.Bytes())
			if "" != this.dir {
				this.data.Write(this.line.Bytes())
			}
			this.line.Reset()
		} else {
			if "" != this.dir {
				this.print2(&this.data, level, file, text, line)
			}
		}
	} else {
		this.print2(&this.line, level, file, text, line)
		if this.list.Len() >= this.recent {
			this.list.Remove(this.list.Front())
		}
		this.list.PushBack(string(this.line.Bytes()))
		if this.std {
			os.Stderr.Write(this.line.Bytes())
		}
		if "" != this.dir {
			this.data.Write(this.line.Bytes())
		}
		this.line.Reset()
	}
	this.Unlock()
}

func (this *log) print2(buf *bytes.Buffer, level Level, file, text *string, line int) {
	this.fmtDateTime.Fmt(time.Now())
	buf.Write(this.fmtDateTime)
	buf.WriteString(*file)
	buf.Write(this.fmtLineNo[:this.fmtLineNo.Fmt(line)])
	buf.Write(levelFmt[level])
	buf.WriteString(*text)
	buf.WriteByte('\n')
}

func (this *log) print3(buf *bytes.Buffer, level Level, fileLine []byte, text *string) {
	this.fmtDateTime.Fmt(time.Now())
	buf.Write(this.fmtDateTime)
	buf.Write(fileLine)
	buf.Write(levelFmt[level])
	buf.WriteString(*text)
	buf.WriteByte('\n')
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
