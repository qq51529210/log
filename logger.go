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

type Logger interface {
	Debug(skip int, text string)
	DebugF(skip int, format string, value ...interface{})
	Warn(skip int, text string)
	WarnF(skip int, format string, value ...interface{})
	Info(skip int, text string)
	InfoF(skip int, format string, value ...interface{})
	Error(skip int, text string)
	ErrorF(skip int, format string, value ...interface{})

	RecoverInside() bool
	RecoverOutside(re interface{}) bool
	Close()
	Recently() []string
}

type logger struct {
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

func (this *logger) Debug(skip int, text string) {
	this.print0(LEVEL_DEBUG, skip+2, &text)
}

func (this *logger) DebugF(skip int, format string, value ...interface{}) {
	this.Debug(skip+1, fmt.Sprintf(format, value...))
}

func (this *logger) Warn(skip int, text string) {
	this.print0(LEVEL_WARN, skip+2, &text)
}

func (this *logger) WarnF(skip int, format string, value ...interface{}) {
	this.Warn(skip+1, fmt.Sprintf(format, value...))
}

func (this *logger) Info(skip int, text string) {
	this.print0(LEVEL_INFO, skip+2, &text)
}

func (this *logger) InfoF(skip int, format string, value ...interface{}) {
	this.Info(skip+1, fmt.Sprintf(format, value...))
}

func (this *logger) Error(skip int, text string) {
	this.print0(LEVEL_ERROR, skip+2, &text)
}

func (this *logger) ErrorF(skip int, format string, value ...interface{}) {
	this.Error(skip+1, fmt.Sprintf(format, value...))
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
		file_line := this.panicFileLine.Find()
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

func (this *logger) print3(buf *bytes.Buffer, level Level, fileLine []byte, text *string) {
	this.fmtDateTime.Fmt(time.Now())
	buf.Write(this.fmtDateTime)
	buf.Write(fileLine)
	buf.Write(levelFmt[level])
	buf.WriteString(*text)
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
