package log

// flume日志，没有实现
type Flume struct {
	rootDir string
}

func (l *Flume) Write(b []byte) (int, error) {
	return 0, nil
}
