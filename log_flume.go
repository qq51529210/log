package log

// flume日志
type LoggerFlume struct {
	rootDir string
}

func (this *LoggerFlume) Write(b []byte) (int, error) {
	return 0, nil
}
