package log

import "io"

// 接管标准库log.Logger的输出
type LoggerStd struct {
	writer   io.Writer
	level    Level
	skip     int
	fileLine FileLine
}

// 某个Package里使用并导出了std.Logger
// 设置std.Logger = NewLoggerStd(server.logger, log.LevelError, 3, log.FileLineFullPath), "", 0)
// 例如http.Server.ErrorLog = NewLoggerStd(server.logger, log.LevelError, 4, log.FileLineFullPath), "", 0)
// 一般Package不会直接使用std.Logger.Output()
// 所以skip一般是3，就可以定位到具体的行
// LoggerStd实现了io.Writer，可以接收std.Logger的输出
func NewLoggerStd(writer io.Writer, level Level, skip int, fileLine FileLine) *LoggerStd {
	return &LoggerStd{writer: writer, level: level, skip: skip, fileLine: fileLine}
}

func (this *LoggerStd) Write(b []byte) (n int, e error) {
	n = len(b)
	if n < 1 {
		return
	}
	if b [n-1] == '\n' {
		n--
	}
	l := logPool.Get().(*Log)
	l.Info.Writer = this.writer
	l.Info.Level = this.level
	l.Info.Skip = this.skip
	l.Info.FileLine = this.fileLine
	// n-1是因为有个\n
	n, e = l.PrintBytes(b[:n])
	logPool.Put(l)
	return
}
