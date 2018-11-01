package log

import (
	"testing"
)

func f0(logger Logger) {
	logger.Print(LevelDebug, 0, "debug")
	f1(logger)
}

func f1(logger Logger) {
	logger.Print(LevelInfo, 1, "info")
	f2(logger)
}

func f2(logger Logger) {
	logger.Print(LevelWarn, 2, "warn")
	f3(logger)
}

func f3(logger Logger) {
	logger.Print(LevelError, 3, "error")
}

func f4() {
	Panic("log panic")
}

func f5(logger Logger) {
	defer func() {
		logger.Recover(recover())
	}()
	f0(logger)
	f4()
}

func f6(logger Logger) {
	defer func() {
		logger.Recover(recover())
	}()
	f0(logger)
	panic("go panic")
}

func TestLog(t *testing.T) {
	l1 := NewStdLogger(LevelDebug,true)
	f5(l1)
	f6(l1)
	l2 := NewStdLogger(LevelDebug,false)
	f5(l2)
	f6(l2)
}
