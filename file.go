package log

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
	"path/filepath"
	"io/ioutil"
	"io"
	"fmt"

	"github.com/qq51529210/common"
)

const (
	MinFileSize = 1024 * 32
	MinDuration = 100 * time.Millisecond
)

var (
	dirFormat  = "20060102"
	fileFormat = "150405.999999999"
)

type FileConfig struct {
	Dir      string `json:"dir"`
	Size     string `json:"size"`
	Day      int    `json:"day"`
	Std      bool   `json:"std"`
	Level    string `json:"level"`
	Stack    string `json:"stack"`
	Duration int    `json:"duration"`
}

func NewFile(cfg *FileConfig) (*File, error) {
	size, e := common.ParseInt(cfg.Size)
	if nil != e {
		return nil, e
	}

	l := &File{
		valid: true,
		dir:   cfg.Dir,
		std:   cfg.Std,
		exit:  make(chan struct{}),
		day:   common.MaxInt(cfg.Day, 1),
		size:  common.MaxInt(int(size), MinFileSize),
		dur:   common.MaxDuration(time.Duration(cfg.Duration)*time.Millisecond, MinDuration),
	}

	if l.dir == "" {
		l.dir = "./"
	}

	switch strings.ToLower(cfg.Stack) {
	case "file":
		l.stack = StackLineFile
	case "path":
		l.stack = StackLinePath
	default:
		l.stack = StackLineNil
	}

	switch strings.ToLower(cfg.Level) {
	case "info":
		l.level = LevelInfo
	case "warn":
		l.level = LevelWarn
	case "error":
		l.level = LevelError
	default:
		l.level = LevelDebug
	}

	l.timer = time.NewTimer(l.dur)
	go l.loop()

	return l, nil
}

type File struct {
	mux   sync.Mutex
	exit  chan struct{}
	dir   string
	size  int
	day   int
	data  bytes.Buffer
	line  bytes.Buffer
	std   bool
	valid bool
	file  *os.File
	level Level
	stack StackLine
	dur   time.Duration
	timer *time.Timer
}

func (this *File) Print(l Level, d int, s string) {
	this.mux.Lock()
	if !this.valid {
		os.Stderr.WriteString("file has been closed")
		this.mux.Unlock()
		return
	}
	if l >= this.level {
		Print(&this.line, l, this.stack, d+1, s)
		if this.std {
			os.Stderr.Write(this.line.Bytes())
		}
		io.Copy(&this.data, &this.line)
		if this.data.Len() >= this.size {
			this.save()
		}
	}
	this.mux.Unlock()
}

func (this *File) Printf(l Level, d int, f string, a ...interface{}) {
	this.Print(l, d+1, fmt.Sprintf(f, a...))
}

func (this *File) Close() error {
	this.mux.Lock()
	if !this.valid {
		this.mux.Unlock()
		return errors.New("file has been closed")
	}
	this.valid = false
	this.mux.Unlock()

	<-this.exit

	return nil
}

func (this *File) save() {
}

func (this *File) new() {
}

func (this *File) close() {
	if nil != this.file {
		this.file.Write(this.data.Bytes())
		this.file.Close()
		this.file = nil
	}
}

func (this *File) loop() {
	defer func() {
		recover()
		this.timer.Stop()
		close(this.exit)
	}()
	for this.valid {
		<-this.timer.C
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
		this.timer.Reset(this.dur)
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
