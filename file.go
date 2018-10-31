package log

import (
	"bytes"
	"os"
	"sync"
	"time"
	"path/filepath"
	"io/ioutil"
	"log"
	"io"
	"fmt"
)

const (
	MinFileSize = 1024 * 32
)

var (
	dirFormat  = "20060102"
	fileFormat = "150405.999999999"
)

func NewFile(dir string, size, day int, std io.Writer) (*File) {
	l := &File{
		dir:   dir,
		valid: true,
		size:  size,
		day:   day,
	}
	l.log = log.New(l, "", 0)
	l.std = std
	if l.dir == "" {
		l.dir = "./"
	}
	if l.size < MinFileSize {
		l.size = MinFileSize
	}
	if l.day < 1 {
		l.day = 1
	}
	go l.saveLoop()
	return l
}

type File struct {
	mux   sync.Mutex
	dir   string
	size  int
	day   int
	data  bytes.Buffer
	valid bool
	file  *os.File
	log   *log.Logger
	std   io.Writer
}

func (this *File) Write(b []byte) (int, error) {
	this.mux.Lock()
	n, e := this.data.Write(b)
	if nil != this.std {
		this.std.Write(b)
	}
	this.mux.Unlock()
	return n, e
}

func (this *File) Print(v ...interface{}) {
	this.log.Output(2, fmt.Sprintln(v...))
}

func (this *File) Printf(v ...interface{}) {
	this.PrintSkip(3,fmt.Sprintln(v...))
}

func (this *File) PrintSkip(n int, v ...interface{}) {
	this.log.Output(n+1, fmt.Sprintln(v...))
}

func (this *File) saveLoop() {
	timer := time.NewTimer(time.Second)
	defer func() {
		recover()
		timer.Stop()
	}()
	for this.valid {
		<-timer.C
		this.mux.Lock()
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
		this.mux.Unlock()
		timer.Reset(time.Second)
	}
}

func (this *File) newFile() {
	now := time.Now()
	dir := filepath.Join(this.dir, now.Format(dirFormat))
	e := os.MkdirAll(dir, os.ModePerm)
	if nil != e {
		os.Stderr.WriteString(e.Error())
		return
	}
	path := filepath.Join(dir, now.Format(fileFormat))
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
			t, e := time.Parse(dirFormat, fs[i].Name())
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

func (this *File) Close() {
	this.mux.Lock()
	if !this.valid {
		this.mux.Unlock()
		return
	}
	this.valid = false
	this.mux.Unlock()

	if nil != this.file {
		this.file.Write(this.data.Bytes())
		this.file.Close()
		this.file = nil
	}
}
