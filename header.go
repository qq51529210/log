package log

import (
	"os"
	"runtime"
	"time"
)

// FormatTime format is "2006-01-02 15:04:05.000000"
func FormatTime(buf *Buffer) {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// Date
	buf.WriteIntLeftAlign(year, 4)
	buf.b = append(buf.b, '-')
	buf.WriteIntRightAlign(int(month), 2)
	buf.b = append(buf.b, '-')
	buf.WriteIntRightAlign(day, 2)
	buf.b = append(buf.b, ' ')
	// Time
	buf.WriteIntRightAlign(hour, 2)
	buf.b = append(buf.b, ':')
	buf.WriteIntRightAlign(minute, 2)
	buf.b = append(buf.b, ':')
	buf.WriteIntRightAlign(second, 2)
	// Nanosecond
	buf.b = append(buf.b, '.')
	buf.WriteInt(t.Nanosecond())
}

type HeaderFunc func(buf *Buffer, depth int)

func DefaultHeader(buf *Buffer, depth int) {
	FormatTime(buf)
}

func FileNameHeader(buf *Buffer, depth int) {
	FormatTime(buf)
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if os.IsPathSeparator(path[i]) {
				path = path[i+1:]
				break
			}
		}
	}
	buf.b = append(buf.b, ' ')
	buf.b = append(buf.b, path...)
	buf.b = append(buf.b, ':')
	buf.WriteInt(line)
}

func FilePathHeader(buf *Buffer, depth int) {
	FormatTime(buf)
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	buf.b = append(buf.b, ' ')
	buf.b = append(buf.b, path...)
	buf.b = append(buf.b, ':')
	buf.WriteInt(line)
}

// // Header 用于格式化日志头
// type Header interface {
// 	// 格式化
// 	Format(buf *Buffer, depth int)
// 	// 格式化，filePath 路径，fileLine 行号
// 	// FormatWith(buf *Buffer, filePath, fileLine string)
// }

// // DeafultHeader 实现 Header 接口
// // 格式 2006-01-02 15:04:05.000000
// type DeafultHeader struct {
// }

// func (th *DeafultHeader) Format(buf *Buffer, depth int) {
// 	FormatTime(buf)
// }

// // func (th *DeafultHeader) FormatWith(buf *Buffer, filePath, fileLine string) {
// // 	FormatTime(buf)
// // }

// // FileNameHeader 实现 Header 接口
// // 格式 2006-01-02 15:04:05.000000 [fileName:fileLine]
// type FileNameHeader struct {
// }

// func (th *FileNameHeader) Format(buf *Buffer, depth int) {
// 	FormatTime(buf)
// }

// // func (th *FileNameHeader) FormatWith(buf *Buffer, filePath, fileLine string) {
// // 	FormatTime(buf)
// // 	i := strings.LastIndexByte(filePath, filepath.Separator)
// // 	if i < 0 {
// // 		buf.b = append(buf.b, ' ')
// // 		buf.b = append(buf.b, filePath...)
// // 		buf.b = append(buf.b, ':')
// // 		buf.b = append(buf.b, fileLine...)
// // 	} else {
// // 		buf.b = append(buf.b, ' ')
// // 		buf.b = append(buf.b, filePath[i+1:]...)
// // 		buf.b = append(buf.b, ':')
// // 		buf.b = append(buf.b, fileLine...)
// // 	}
// // }

// // FilePathHeader 实现 Header 接口
// // 格式 2006-01-02 15:04:05.000000 [filePath:fileLine]
// type FilePathHeader struct {
// }

// func (th *FilePathHeader) Format(buf *Buffer, depth int) {
// 	FormatTime(buf)
// }

// // func (th *FilePathHeader) FormatWith(buf *Buffer, filePath, fileLine string) {
// // 	FormatTime(buf)
// // 	buf.b = append(buf.b, ' ')
// // 	buf.b = append(buf.b, filePath...)
// // 	buf.b = append(buf.b, ':')
// // 	buf.b = append(buf.b, fileLine...)
// // }
