package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	myLogger  = NewLogger(os.Stderr, 1, new(FilePathStackHeader))
	logFormat = "text: %s"
	logText   = "test"
)

func Test_Logger(t *testing.T) {
	test_Debug()
	test_Info()
	test_Warn()
	test_Error()
}

func test_Debug() {
	Debug("Debug:", logText)
	Debugf("Debugf: %s", logText)
	DepthDebug(1, "DepthDebug:", logText)
	DepthDebugf(1, "DepthDebugf: %s", logText)
}

func test_Info() {
	Info("Info:", logText)
	Infof("Infof: %s", logText)
	DepthInfo(1, "DepthInfo:", logText)
	DepthInfof(1, "DepthInfof: %s", logText)
}

func test_Warn() {
	Warn("Warn:", logText)
	Warnf("Warnf: %s", logText)
	DepthWarn(1, "DepthWarn:", logText)
	DepthWarnf(1, "DepthWarnf: %s", logText)
}

func test_Error() {
	Error("Error:", logText)
	Errorf("Errorf: %s", logText)
	DepthError(1, "DepthError:", logText)
	DepthErrorf(1, "DepthErrorf: %s", logText)
}

func Test_Recover(t *testing.T) {
	defer func() {
		Recover(recover())
	}()

	panic("test recover")
}

func Benchmark_My_Logger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	myLogger.SetWriter(&output)
	for i := 0; i < b.N; i++ {
		myLogger.Debug(logText)
		output.Reset()
	}
}

func Benchmark_Std_Logger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Println(logText)
		output.Reset()
	}
}

func Benchmark_My_Logger_f(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	myLogger.SetWriter(&output)
	for i := 0; i < b.N; i++ {
		myLogger.Debugf(logFormat, logText)
		output.Reset()
	}
}

func Benchmark_Std_Logger_f(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Printf(logFormat, logText)
		output.Reset()
	}
}
