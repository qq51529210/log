package log

import (
	"path/filepath"
	"runtime"
	"time"
)

type HeaderFormaterType string

const (
	FileNameStackHeaderFormater HeaderFormaterType = "fileNameStack"
	FilePathStackHeaderFormater HeaderFormaterType = "filePathStack"
)

// HeaderFormater 格式化日志头
type HeaderFormater interface {
	// FormatWith 使用 level ，path ，line 来格式化 log 。
	// path 和 line 是调用堆栈的文件路径和行数。
	FormatWith(log *Log, trackId string, level Level, path, line string)
	// Format 使用 level ，depth 来格式化 log 。
	// depth 是 runtime.Caller() 的参数。
	Format(log *Log, trackId string, level Level, depth int)
}

// NewHeaderFormater 返回不同的 HeaderFormater 。
// fileNameStack 返回的是 addId level datetime fileName:line log 的格式，
// filePathStack 返回的是 addId level datetime filePath:line log 的格式。
func NewHeaderFormater(headerFormaterType HeaderFormaterType, appId string) HeaderFormater {
	switch headerFormaterType {
	case FileNameStackHeaderFormater:
		return &fileNameStackHeader{defaultHeader: defaultHeader{appId: appId}}
	case FilePathStackHeaderFormater:
		return &filePathStackHeader{defaultHeader: defaultHeader{appId: appId}}
	default:
		return &defaultHeader{appId: appId}
	}
}

// defaultHeader format: AppId Level Date Time
type defaultHeader struct {
	appId string
}

func (h *defaultHeader) Format(log *Log, trackId string, level Level, depth int) {
	// level
	log.line = append(log.line, byte(level))
	log.line = append(log.line, ' ')
	// app id
	if h.appId != "" {
		log.line = append(log.line, h.appId...)
		log.line = append(log.line, ' ')
	}
	// track id
	if trackId != "" {
		log.line = append(log.line, trackId...)
		log.line = append(log.line, ' ')
	}
	// time
	FormatTime(log)
}

func (h *defaultHeader) FormatWith(log *Log, trackId string, level Level, path, line string) {
	// level
	log.line = append(log.line, byte(level))
	log.line = append(log.line, ' ')
	// app id
	if h.appId != "" {
		log.line = append(log.line, h.appId...)
		log.line = append(log.line, ' ')
	}
	// track id
	if trackId != "" {
		log.line = append(log.line, trackId...)
		log.line = append(log.line, ' ')
	}
	// time
	FormatTime(log)
}

// fileNameStackHeader format is "AppId Level Date Time StackFileName:CodeLine"
type fileNameStackHeader struct {
	defaultHeader
}

func (h *fileNameStackHeader) Format(log *Log, trackId string, level Level, depth int) {
	h.defaultHeader.Format(log, trackId, level, depth)
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

func (h *fileNameStackHeader) FormatWith(log *Log, trackId string, level Level, path, line string) {
	h.defaultHeader.FormatWith(log, trackId, level, path, line)
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

// fileNameStackHeader format is "AppId Level Date Time StackFilePath:CodeLine"
type filePathStackHeader struct {
	defaultHeader
}

func (h *filePathStackHeader) Format(log *Log, trackId string, level Level, depth int) {
	h.defaultHeader.Format(log, trackId, level, depth)
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

func (h *filePathStackHeader) FormatWith(log *Log, trackId string, level Level, path, line string) {
	h.defaultHeader.FormatWith(log, trackId, level, path, line)
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
