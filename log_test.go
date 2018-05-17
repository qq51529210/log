package log

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutput(t *testing.T) {
	dir := filepath.Join(os.Getenv("GOPATH"), "src/github.com/qq51529210/log/test")
	l := Open(dir, 64*1024, 3, 100, true)

	l.Debug("Debug")
	l.Warn("Warn")
	l.Info("Info")
	l.Error("Error")

	l.DebugF("DebugF %v", 1)
	l.WarnF("WarnF %v", 2)
	l.InfoF("InfoF %v", 3)
	l.ErrorF("ErrorF %v", 4)

	l.DebugSkip(0, "DebugSkip")
	l.WarnSkip(0, "WarnSkip")
	l.InfoSkip(0, "InfoSkip")
	l.ErrorSkip(0, "ErrorSkip")

	l.DebugSkipF(0, "DebugSkipF")
	l.WarnSkipF(0, "WarnSkipF")
	l.InfoSkipF(0, "InfoSkipF")
	l.ErrorSkipF(0, "ErrorSkipF")

	F1(l)
	F2(l)
	F3(l)
	F4(l)

	l.Close()

	os.RemoveAll(dir)
}

func F1(l Logger) {
	defer func() {
		l.RecoverOutside(recover())
	}()

	l.Debug("Debug")
	l.Warn("Warn")
	l.Info("Info")
	l.Error("Error")

	Panic("Panic")
}

func F2(l Logger) {
	defer l.RecoverInside()

	l.DebugF("DebugF %v", 1)
	l.WarnF("WarnF %v", 2)
	l.InfoF("InfoF %v", 3)
	l.ErrorF("ErrorF %v", 4)

	PanicF("PanicF %v", 5)
}

func F3(l Logger) {
	defer func() {
		l.RecoverOutside(recover())
	}()

	l.DebugSkip(1, "DebugSkip")
	l.WarnSkip(1, "WarnSkip")
	l.InfoSkip(1, "InfoSkip")
	l.ErrorSkip(1, "ErrorSkip")

	PanicSkip(1, "PanicSkip")
}

func F4(l Logger) {
	defer l.RecoverInside()

	l.DebugSkipF(1, "DebugSkipF %v", 1)
	l.WarnSkipF(1, "WarnSkipF %v", 2)
	l.InfoSkipF(1, "InfoSkipF %v", 3)
	l.ErrorSkipF(1, "ErrorSkipF %v", 4)

	PanicSkipF(1, "PanicSkipF %v", 5)
}
