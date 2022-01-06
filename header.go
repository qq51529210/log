package log

import (
	"path/filepath"
	"runtime"
	"time"
)

type Header interface {
	FormatWith(log *Log, path, line string)
	Format(log *Log, depth int)
}

type NullStackHeader struct {
}

func (hd *NullStackHeader) Format(log *Log, depth int) {
	FormatTime(log)
}

func (hd *NullStackHeader) FormatWith(log *Log, path, line string) {
	FormatTime(log)
}

type FileNameStackHeader struct {
}

func (hd *FileNameStackHeader) Format(log *Log, depth int) {
	FormatTime(log)
	log.line = append(log.line, ' ')
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if path[i] == filepath.Separator {
				path = path[i+1:]
				break
			}
		}
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
}

func (hd *FileNameStackHeader) FormatWith(log *Log, path, line string) {
	FormatTime(log)
	log.line = append(log.line, ' ')
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == filepath.Separator {
			path = path[i+1:]
			break
		}
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteString(line)
}

type FilePathStackHeader struct {
}

func (hd *FilePathStackHeader) Format(log *Log, depth int) {
	FormatTime(log)
	log.line = append(log.line, ' ')
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
}

func (hd *FilePathStackHeader) FormatWith(log *Log, path, line string) {
	FormatTime(log)
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteString(line)
}

func FormatTime(log *Log) {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// Date
	log.WriteLeftAlignInt(year, 4)
	log.line = append(log.line, '-')
	log.WriteRightAlignInt(int(month), 2)
	log.line = append(log.line, '-')
	log.WriteRightAlignInt(day, 2)
	log.line = append(log.line, ' ')
	// Time
	log.WriteRightAlignInt(hour, 2)
	log.line = append(log.line, ':')
	log.WriteRightAlignInt(minute, 2)
	log.line = append(log.line, ':')
	log.WriteRightAlignInt(second, 2)
	// Nanosecond
	log.line = append(log.line, '.')
	log.WriteLeftAlignInt(t.Nanosecond(), 6)
}
