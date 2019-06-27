package log

import "io"

type LoggerStd struct {
	writer   io.Writer
	level    Level
	skip     int
	fileLine FileLine
}

// skip 一般是3
func NewLoggerStd(writer io.Writer, level Level, skip int, fileLine FileLine) *LoggerStd {
	return &LoggerStd{
		writer:   writer,
		level:    level,
		skip:     skip,
		fileLine: fileLine,
	}
}

func (this *LoggerStd) Write(b []byte) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.PrintBytes(this.writer, this.level, this.skip, this.fileLine, b)
	logPool.Put(l)
	return n, e
}
