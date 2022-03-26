package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	myLogger  = NewLogger(os.Stderr, NewHeaderFormater(FilePathStackHeaderFormater, "appId"))
	logFormat = "text: %s"
	logText   = "test"
)

func init() {
	SetHeaderFormater(NewHeaderFormater(FilePathStackHeaderFormater, "test-app"))
}

func Test_Logger(t *testing.T) {
	test_Debug()
	test_Info()
	test_Warn()
	test_Error()
}

func test_Debug() {
	Debug("Debug:", logText)
	Debugf("Debugf: %s", logText)
	DepthDebug("track-id-1", 1, "DepthDebug:", logText)
	DepthDebugf("track-id-2", 1, "DepthDebugf: %s", logText)
}

func test_Info() {
	Info("Info:", logText)
	Infof("Infof: %s", logText)
	DepthInfo("track-id-1", 1, "DepthInfo:", logText)
	DepthInfof("track-id-1", 1, "DepthInfof: %s", logText)
}

func test_Warn() {
	Warn("Warn:", logText)
	Warnf("Warnf: %s", logText)
	DepthWarn("track-id-11", 1, "DepthWarn:", logText)
	DepthWarnf("track-id-12", 1, "DepthWarnf: %s", logText)
}

func test_Error() {
	Error("Error:", logText)
	Errorf("Errorf: %s", logText)
	DepthError("track-id-111", 1, "DepthError:", logText)
	DepthErrorf("track-id-112", 1, "DepthErrorf: %s", logText)
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
	myLogger.SetOutput(&output)
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
	myLogger.Info("benchmark my logger")
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	myLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		myLogger.Debugf(logFormat, logText)
		output.Reset()
	}
}

func Benchmark_Std_Logger_f(b *testing.B) {
	stdLogger.Println("benchmark std logger")
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Printf(logFormat, logText)
		output.Reset()
	}
}
