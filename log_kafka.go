package log

// kafka日志，没有实现
type LoggerKafka struct {
	rootDir string
}

func (this *LoggerKafka) Write(b []byte) (int, error) {
	return 0, nil
}
