package log

import (
	"path/filepath"
	"runtime"
	"time"
)

// HeaderFormaterType 表示格式化日志头的类型
type HeaderFormaterType string

const (
	// DefaultHeaderFormater 表示默认的格式化头
	DefaultHeaderFormater HeaderFormaterType = "default"
	// FileNameStackHeaderFormater 表示打印堆栈信息只打印文件名
	FileNameStackHeaderFormater HeaderFormaterType = "fileNameStack"
	// FilePathStackHeaderFormater 表示打印堆栈信息只打印文件完整路径
	FilePathStackHeaderFormater HeaderFormaterType = "filePathStack"
)

// HeaderFormater 格式化日志头
type HeaderFormater interface {
	// FormatWith 使用 level ，path ，line 来格式化 log 。
	// path 和 line 是调用堆栈的文件路径和行数。
	FormatWith(log *Log, traceID string, level Level, path, line string)
	// Format 使用 level ，depth 来格式化 log 。
	// depth 是 runtime.Caller() 的参数。
	Format(log *Log, traceID string, level Level, depth int)
}

// NewHeaderFormater 返回不同的 HeaderFormater 。
// fileNameStack 返回的是 appID level datetime fileName:line log 的格式，
// filePathStack 返回的是 appID level datetime filePath:line log 的格式。
func NewHeaderFormater(headerFormaterType HeaderFormaterType, appID string) HeaderFormater {
	switch headerFormaterType {
	case FileNameStackHeaderFormater:
		return &fileNameStackHeader{defaultHeader: defaultHeader{appID: appID}}
	case FilePathStackHeaderFormater:
		return &filePathStackHeader{defaultHeader: defaultHeader{appID: appID}}
	default:
		return &defaultHeader{appID: appID}
	}
}

// defaultHeader format: AppID Level Date Time
type defaultHeader struct {
	appID string
}

func (h *defaultHeader) Format(log *Log, traceID string, level Level, depth int) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
}

func (h *defaultHeader) FormatWith(log *Log, traceID string, level Level, path, line string) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
}

// fileNameStackHeader format is "AppID Level Date Time StackFileName:CodeLine"
type fileNameStackHeader struct {
	defaultHeader
}

func (h *fileNameStackHeader) Format(log *Log, traceID string, level Level, depth int) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// path line
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
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
	log.line = append(log.line, ':')
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
}

func (h *fileNameStackHeader) FormatWith(log *Log, traceID string, level Level, path, line string) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// path line
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == filepath.Separator {
			path = path[i+1:]
			break
		}
	}
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteString(line)
	log.line = append(log.line, ':')
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
}

// fileNameStackHeader format is "AppID Level Date Time StackFilePath:CodeLine"
type filePathStackHeader struct {
	defaultHeader
}

func (h *filePathStackHeader) Format(log *Log, traceID string, level Level, depth int) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// path line
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteInt(line)
	log.line = append(log.line, ':')
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
}

func (h *filePathStackHeader) FormatWith(log *Log, traceID string, level Level, path, line string) {
	// level
	log.line = append(log.line, byte(level))
	// app id
	if h.appID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, h.appID...)
	}
	// time
	log.line = append(log.line, ' ')
	FormatTime(log)
	// path line
	log.line = append(log.line, ' ')
	log.WriteString(path)
	log.line = append(log.line, ':')
	log.WriteString(line)
	log.line = append(log.line, ':')
	// trace id
	if traceID != "" {
		log.line = append(log.line, ' ')
		log.line = append(log.line, traceID...)
	}
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
