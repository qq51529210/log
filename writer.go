package log

import (
	"errors"
)

var (
	errEmptyBuffer = errors.New("empty buffer")
)

// 实现了io.Writer，可以接管标准库log.Logger的输出
// 某个Package里使用并导出了std.Logger
// 设置std.Logger = new(Writer)
// 例如http.Server.ErrorLog = new(Writer)
// 一般Package不会直接使用std.Logger.Output()
// 所以skip一般是3，就可以定位到具体的行
type Writer struct {
	Skip int
	Level
}

func (w *Writer) Write(b []byte) (int, error) {
	n := len(b)
	if n < 1 {
		return 0, errEmptyBuffer
	}
	if b[n-1] == '\n' {
		n--
	}
	l := logPool.Get().(*Logger)
	_, err := l.PrintBytes(defaultWriter, w.Level, defaultStack, w.Skip, b[:n])
	logPool.Put(l)
	return n, err
}
