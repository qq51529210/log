package log

// kafka日志
type LoggerKafka struct {
	rootDir string
}

func (this *LoggerKafka) Write(b []byte) (int, error) {
	return 0, nil
}
