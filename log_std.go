package log

import "io"

type LoggerStd struct {
	writer   io.Writer
	level    Level
	skip     int
	fileLine FileLine
}

// 某个Package里使用并导出了std.Logger
// 设置std.Logger = NewLoggerStd(server.logger, log.LevelError, 3, log.FileLineFullPath), "", 0)
// 一般Package不会直接使用std.Logger.Output()
// 所以skip一般是3，就可以定位到具体的行
// LoggerStd实现了io.Writer，可以接收std.Logger的输出
func NewLoggerStd(writer io.Writer, level Level, skip int, fileLine FileLine) *LoggerStd {
	return &LoggerStd{writer: writer, level: level, skip: skip, fileLine: fileLine}
}

func (this *LoggerStd) Write(b []byte) (int, error) {
	l := logPool.Get().(*Log)
	n, e := l.PrintBytes(this.writer, this.level, this.skip, this.fileLine, b)
	logPool.Put(l)
	return n, e
}
