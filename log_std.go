package log

import "io"

type LoggerStd struct {
	writer   io.Writer
	level    Level
	skip     int
	fileLine FileLine
}

func NewLoggerStd(writer io.Writer, level Level, fileLine FileLine, skip int) *LoggerStd {
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
