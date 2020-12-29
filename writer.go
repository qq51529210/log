package log

import (
	"errors"
	"io"
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
	Out io.Writer
}

func (w *Writer) Write(b []byte) (int, error) {
	n := len(b)
	if n < 1 {
		return 0, errEmptyBuffer
	}
	if b[n-1] == '\n' {
		n--
	}
	l := logPool.Get().(*Log)
	n, err := l.PrintBytes(w.Out, w.Level, w.Skip, b)
	logPool.Put(l)
	return n, err
}
