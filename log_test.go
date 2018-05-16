package log

import (
	"errors"
	"sync"
	"testing"
)

func TestOutput(t *testing.T) {
	l := Open("", 0, 0, 0, true)
	l.Debug(0, "Debug")
	l.Warn(0, "Warn")
	l.Info(0, "Info")
	l.Error(0, "Error")

	l.DebugF(0, "Debug %v", 1)
	l.WarnF(0, "Warn %v", 2)
	l.InfoF(0, "Info %v", 3)
	l.ErrorF(0, "Error %v", 4)

	l.Close()
}
