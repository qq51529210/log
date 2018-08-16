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
	"strings"
	"strconv"
)

const (
	_LEVEL_      Level = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	_LEVEL_PANIC
	_LEVEL_STACK
)

const (
	kb = 1024
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
)

type Level int

var (
	date_dir_fmt  = "20060102"
	time_file_fmt = "150405.999999999"
	levelFmt      = [][]byte{
		[]byte(string(" [NULL] ")),
		[]byte(string(" [DEBUG] ")),
		[]byte(string(" [INFO] ")),
		[]byte(string(" [WARN] ")),
		[]byte(string(" [ERROR] ")),
		[]byte(string(" [PANIC] ")),
		[]byte(string(" [STACK] ")),
	}
	unknownFileLine = [][]byte{[]byte("???:-1")}
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

type Config struct {
	Level    string `json:"level"`
	Dir      string `json:"dir"`
	Size     string `json:"size"`
	Day      int    `json:"day"`
	Resent   int    `json:"resent"`
	STD      bool   `json:"std"`
	FullPath bool   `json:"full_path"`
}

type panicInfo struct {
	file  string
	line  int
	value interface{}
	time  time.Time
}

func (this *panicInfo) String() string {
	dt := newFmtDateTime()
	dt.Fmt(this.time)
	s := string(dt)
	s += this.file
	ln := newFmtLineNo()
	s += string(ln[:ln.Fmt(this.line)])
	s += string(levelFmt[_LEVEL_PANIC])
	s += fmt.Sprintf("%v\n", this.value)
	return s
}

type panicFileLine []byte

func newPanicFileLine() panicFileLine {
	return make([]byte, 32)
}

func (this panicFileLine) Find(fullPath bool) [][]byte {
	lines := make([][]byte, 0)
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
		// read line
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
				// file line
				if line[0] == '\t' {
					j := 1
					for ; j < len(line); j++ {
						if line[j] == ' ' {
							break
						}
					}
					if !fullPath {
						for k := j - 1; k >= 1; k-- {
							if line[k] == '/' {
								l := make([]byte, j-k)
								copy(l, line[k+1:j])
								lines = append(lines, l)
								break
							}
						}
					} else {
						l := make([]byte, j-1)
						copy(l, line[1:j])
						lines = append(lines, l)
					}
				}
			}
			return lines
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

func parseBytes(s string) (uint64, error) {
	s = strings.ToLower(s)
	s = strings.Replace(s, "b", "", -1)
	if strings.Contains(s, "t") {
		s = strings.Replace(s, "t", "", -1)
		v, e := strconv.ParseFloat(s, 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * tb), nil
	}
	if strings.Contains(s, "g") {
		s = strings.Replace(s, "g", "", -1)
		v, e := strconv.ParseFloat(s, 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * gb), nil
	}
	if strings.Contains(s, "m") {
		s = strings.Replace(s, "m", "", -1)
		v, e := strconv.ParseFloat(s, 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * mb), nil
	}
	if strings.Contains(s, "k") {
		s = strings.Replace(s, "k", "", -1)
		v, e := strconv.ParseFloat(s, 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * kb), nil
	}
	v, e := strconv.ParseFloat(s, 64)
	if nil != e {
		return 0, e
	}
	return uint64(v), nil
}

func OpenWith(cfg *Config) Logger {
	size, e := parseBytes(cfg.Size)
	if nil != e {
		panic(e)
	}
	return Open(cfg.Dir, cfg.Level, int(size), cfg.Day, cfg.Resent, cfg.STD, cfg.FullPath)
}

func Open(dir, level string, size, day, recent int, std, fullPath bool) Logger {
	l := &logger{
		dir:           dir,
		valid:         true,
		fullPath:      fullPath,
		std:           std,
		recent:        recent,
		stack:         make([]byte, 4096),
		fmtDateTime:   newFmtDateTime(),
		fmtLineNo:     newFmtLineNo(),
		panicFileLine: newPanicFileLine(),
	}
	l.SetLevel(level)
	l.SetSize(size)
	l.SetDay(day)
	if l.dir != "" {
		l.Add(1)
		go l.syncLoop()
	}
	return l
}

func PrintF(level Level, format string, value ...interface{}) {
	PrintSkip(level, 1, fmt.Sprintf(format, value...))
}

func PrintSkipF(level Level, skip int, format string, value ...interface{}) {
	PrintSkip(level, skip+1, fmt.Sprintf(format, value...))
}

func Print(level Level, text string) {
	_log.print0(level, 2, &text)
}

func PrintSkip(level Level, skip int, text string) {
	_log.print0(level, skip+2, &text)
}

func CheckError(e error) {
	if nil == e {
		return
	}
	PanicSkip(1, e)
}

func Panic(value interface{}) {
	PanicSkip(1, value)
}

func PanicSkip(skip int, value interface{}) {
	o := false
	info := &panicInfo{time: time.Now(), value: value}
	_, info.file, info.line, o = runtime.Caller(skip + 1)
	if !o {
		info.file = "???"
		info.line = -1
	}
	panic(info)
}

func PanicF(format string, value ... interface{}) {
	PanicSkip(1, fmt.Sprintf(format, value...))
}

func PanicSkipF(skip int, format string, value ... interface{}) {
	PanicSkip(skip+1, fmt.Sprintf(format, value...))
}

type Logger interface {
	Debug(text string)
	DebugF(format string, value ...interface{})
	Warn(text string)
	WarnF(format string, value ...interface{})
	Info(text string)
	InfoF(format string, value ...interface{})
	Error(text string)
	ErrorF(format string, value ...interface{})

	DebugSkip(skip int, text string)
	DebugSkipF(skip int, format string, value ...interface{})
	WarnSkip(skip int, text string)
	WarnSkipF(skip int, format string, value ...interface{})
	InfoSkip(skip int, text string)
	InfoSkipF(skip int, format string, value ...interface{})
	ErrorSkip(skip int, text string)
	ErrorSkipF(skip int, format string, value ...interface{})

	SetSize(int)
	SetDay(int)
	SetRecent(int)
	SetLevel(string)
	EnableStdOutput(bool)
	RecoverInside() bool
	RecoverOutside(re interface{}) bool
	Close()
	Recently() []string
}

type logger struct {
	sync.RWMutex
	sync.WaitGroup
	level    Level
	dir      string
	size     int
	day      int
	recent   int
	std      bool
	data     bytes.Buffer
	valid    bool
	fullPath bool
	stack    []byte
	line     bytes.Buffer
	file     *os.File
	list     list.List
	fmtDateTime
	fmtLineNo
	panicFileLine
}

func (this *logger) SetLevel(level string) {
	switch level {
	case "debug":
		this.level = LEVEL_DEBUG
	case "info":
		this.level = LEVEL_INFO
	case "warn":
		this.level = LEVEL_WARN
	case "error":
		this.level = LEVEL_ERROR
	case "panic":
		this.level = _LEVEL_PANIC
	default:
		this.level = _LEVEL_
	}
}

func (this *logger) SetSize(n int) {
	if n < 1 {
		n = 1024 * 1024
	}
	this.size = n
}

func (this *logger) SetDay(n int) {
	if n < 1 {
		n = 7
	}
	this.day = n
}

func (this *logger) SetRecent(n int) {
	this.recent = n
}

func (this *logger) EnableStdOutput(enable bool) {
	this.std = enable
}

func (this *logger) Debug(text string) {
	if this.level <= LEVEL_DEBUG {
		this.print0(LEVEL_DEBUG, 2, &text)
	}
}

func (this *logger) DebugF(format string, value ...interface{}) {
	if this.level <= LEVEL_DEBUG {
		this.DebugSkip(1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) Warn(text string) {
	if this.level <= LEVEL_WARN {
		this.print0(LEVEL_WARN, 2, &text)
	}
}

func (this *logger) WarnF(format string, value ...interface{}) {
	if this.level <= LEVEL_WARN {
		this.WarnSkip(1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) Info(text string) {
	if this.level <= LEVEL_INFO {
		this.print0(LEVEL_INFO, 2, &text)
	}
}

func (this *logger) InfoF(format string, value ...interface{}) {
	if this.level <= LEVEL_INFO {
		this.InfoSkip(1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) Error(text string) {
	if this.level <= LEVEL_ERROR {
		this.print0(LEVEL_ERROR, 2, &text)
	}
}

func (this *logger) ErrorF(format string, value ...interface{}) {
	if this.level <= LEVEL_ERROR {
		this.ErrorSkip(1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) DebugSkip(skip int, text string) {
	if this.level <= LEVEL_DEBUG {
		this.print0(LEVEL_DEBUG, skip+2, &text)
	}
}

func (this *logger) DebugSkipF(skip int, format string, value ...interface{}) {
	if this.level <= LEVEL_DEBUG {
		this.DebugSkip(skip+1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) WarnSkip(skip int, text string) {
	if this.level <= LEVEL_WARN {
		this.print0(LEVEL_WARN, skip+2, &text)
	}
}

func (this *logger) WarnSkipF(skip int, format string, value ...interface{}) {
	if this.level <= LEVEL_WARN {
		this.WarnSkip(skip+1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) InfoSkip(skip int, text string) {
	if this.level <= LEVEL_INFO {
		this.print0(LEVEL_INFO, skip+2, &text)
	}
}

func (this *logger) InfoSkipF(skip int, format string, value ...interface{}) {
	if this.level <= LEVEL_INFO {
		this.InfoSkip(skip+1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) ErrorSkip(skip int, text string) {
	if this.level <= LEVEL_ERROR {
		this.print0(LEVEL_ERROR, skip+2, &text)
	}
}

func (this *logger) ErrorSkipF(skip int, format string, value ...interface{}) {
	if this.level <= LEVEL_ERROR {
		this.ErrorSkip(skip+1, fmt.Sprintf(format, value...))
	}
}

func (this *logger) RecoverInside() bool {
	return this.RecoverOutside(recover())
}

func (this *logger) RecoverOutside(re interface{}) bool {
	if nil == re {
		return false
	}
	switch re.(type) {
	case *panicInfo:
		info := re.(*panicInfo)
		text := fmt.Sprint(info.value)
		this.print1(_LEVEL_PANIC, &info.file, &text, info.line)
	default:
		text := fmt.Sprint(re)
		file_line := this.panicFileLine.Find(this.fullPath)
		this.Lock()
		switch this.recent {
		case 0:
			if this.std {
				this.print3(&this.line, _LEVEL_PANIC, file_line, &text)
				os.Stderr.Write(this.line.Bytes())
				if "" != this.dir {
					this.data.Write(this.line.Bytes())
				}
				this.line.Reset()
			} else {
				if "" != this.dir {
					this.print3(&this.data, _LEVEL_PANIC, file_line, &text)
				}
			}
		default:
			this.print3(&this.line, _LEVEL_PANIC, file_line, &text)
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
	return true
}

func (this *logger) Close() {
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

func (this *logger) Recently() []string {
	recent := make([]string, 0)
	this.RLock()
	for ele := this.list.Front(); nil != ele; ele = ele.Next() {
		recent = append(recent, ele.Value.(string))
	}
	this.RUnlock()
	return recent
}

func (this *logger) print0(level Level, skip int, text *string) {
	_, f, l, o := runtime.Caller(skip)
	if !o {
		f = "???"
		l = -1
	}
	if !this.fullPath {
		for i := len(f) - 1; i >= 0; i-- {
			if f[i] == '/' {
				f = f[i+1:]
				break
			}
		}
	}
	this.print1(level, &f, text, l)
}

func (this *logger) print1(level Level, file, text *string, line int) {
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

func (this *logger) print2(buf *bytes.Buffer, level Level, file, text *string, line int) {
	this.fmtDateTime.Fmt(time.Now())
	buf.Write(this.fmtDateTime)
	buf.WriteString(*file)
	buf.Write(this.fmtLineNo[:this.fmtLineNo.Fmt(line)])
	buf.Write(levelFmt[level])
	buf.WriteString(*text)
	buf.WriteByte('\n')
}

func (this *logger) print3(buf *bytes.Buffer, level Level, fileLine [][]byte, text *string) {
	this.fmtDateTime.Fmt(time.Now())
	buf.Write(this.fmtDateTime)
	buf.Write(fileLine[0])
	buf.Write(levelFmt[level])
	buf.WriteString(*text)
	buf.Write(levelFmt[_LEVEL_STACK])
	for i := 1; i < len(fileLine); i++ {
		buf.Write(fileLine[i])
		buf.WriteByte(' ')
	}
	buf.WriteByte('\n')
}

func (this *logger) syncLoop() {
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

func (this *logger) newFile() {
	now := time.Now()
	dir := filepath.Join(this.dir, now.Format(date_dir_fmt))
	e := os.MkdirAll(dir, os.ModePerm)
	if nil != e {
		os.Stderr.WriteString(e.Error())
		return
	}
	path := filepath.Join(dir, now.Format(time_file_fmt))
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
