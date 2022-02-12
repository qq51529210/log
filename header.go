package log

import (
	"path/filepath"
	"runtime"
	"time"
)

type Header interface {
	FormatWith(log *Log, level Level, path, line string)
	Format(log *Log, level Level, depth int)
}

// DefaultHeader format: AppId Level Date Time
type DefaultHeader struct {
	id string
}

func (h *DefaultHeader) Format(log *Log, level Level, depth int) {
	// app id
	log.line = append(log.line, h.id...)
	log.line = append(log.line, ' ')
	// level
	log.line = append(log.line, byte(level))
	log.line = append(log.line, ' ')
	// time
	FormatTime(log)
}

func (h *DefaultHeader) FormatWith(log *Log, level Level, path, line string) {
	// app id
	log.line = append(log.line, h.id...)
	log.line = append(log.line, ' ')
	// level
	log.line = append(log.line, byte(level))
	log.line = append(log.line, ' ')
	// time
	FormatTime(log)
}

// FileNameStackHeader format is "AppId Level Date Time StackFileName:CodeLine"
type FileNameStackHeader struct {
	DefaultHeader
}

func (h *FileNameStackHeader) Format(log *Log, level Level, depth int) {
	h.DefaultHeader.Format(log, level, depth)
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

func (h *FileNameStackHeader) FormatWith(log *Log, level Level, path, line string) {
	h.DefaultHeader.FormatWith(log, level, path, line)
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

// FileNameStackHeader format is "AppId Level Date Time StackFilePath:CodeLine"
type FilePathStackHeader struct {
	DefaultHeader
}

func (h *FilePathStackHeader) Format(log *Log, level Level, depth int) {
	h.DefaultHeader.Format(log, level, depth)
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

func (h *FilePathStackHeader) FormatWith(log *Log, level Level, path, line string) {
	h.DefaultHeader.FormatWith(log, level, path, line)
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteString(line)
}

// FormatTime format is "2006-01-02 15:04:05.000000"
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
