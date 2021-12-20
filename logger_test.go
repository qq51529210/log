package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	myLogger  = NewLogger(os.Stderr, 1, new(CallStackFilePathHeader))
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
	LevelDebug(0, "LevelDebug:", logText)
	LevelDebugf(-1, "LevelDebugf: %s", logText)
	LevelDepthDebug(0, 0, "LevelDepthDebug:", logText)
	LevelDepthDebug(-1, 1, "LevelDepthDebug:", logText)
	LevelDepthDebugf(0, 0, "LevelDepthDebugf: %s", logText)
	LevelDepthDebugf(-1, 1, "LevelDepthDebugf: %s", logText)
}

func test_Info() {
	Info("Info:", logText)
	Infof("Infof: %s", logText)
	DepthInfo(1, "DepthInfo:", logText)
	DepthInfof(1, "DepthInfof: %s", logText)
	LevelInfo(0, "LevelInfo:", logText)
	LevelInfof(-1, "LevelInfof: %s", logText)
	LevelDepthInfo(0, 0, "LevelDepthInfo:", logText)
	LevelDepthInfo(-1, 1, "LevelDepthInfo:", logText)
	LevelDepthInfof(0, 0, "LevelDepthInfof: %s", logText)
	LevelDepthInfof(-1, 1, "LevelDepthInfof: %s", logText)
}

func test_Warn() {
	Warn("Warn:", logText)
	Warnf("Warnf: %s", logText)
	DepthWarn(1, "DepthWarn:", logText)
	DepthWarnf(1, "DepthWarnf: %s", logText)
	LevelWarn(0, "LevelWarn:", logText)
	LevelWarnf(-1, "LevelWarnf: %s", logText)
	LevelDepthWarn(0, 0, "LevelDepthWarn:", logText)
	LevelDepthWarn(-1, 1, "LevelDepthWarn:", logText)
	LevelDepthWarnf(0, 0, "LevelDepthWarnf: %s", logText)
	LevelDepthWarnf(-1, 1, "LevelDepthWarnf: %s", logText)
}

func test_Error() {
	Error("Error:", logText)
	Errorf("Errorf: %s", logText)
	DepthError(1, "DepthError:", logText)
	DepthErrorf(1, "DepthErrorf: %s", logText)
	LevelError(0, "LevelError:", logText)
	LevelErrorf(-1, "LevelErrorf: %s", logText)
	LevelDepthError(0, 0, "LevelDepthError:", logText)
	LevelDepthError(-1, 1, "LevelDepthError:", logText)
	LevelDepthErrorf(0, 0, "LevelDepthErrorf: %s", logText)
	LevelDepthErrorf(-1, 1, "LevelDepthErrorf: %s", logText)
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
