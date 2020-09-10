package log

// kafka日志，没有实现
type Kafka struct {
}

func (k *Kafka) Write(b []byte) (int, error) {
	return 0, nil
}
