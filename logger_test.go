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
	DebugDepth(0, "DebugDepth:", logText)
	DebugfDepth(1, "DebugDepthf: %s", logText)
	DebugTrack("track-id-1", "Debug:", logText)
	DebugfTrack("track-id-2", "Debugf: %s", logText)
	DebugDepthTrack("track-id-3", 0, "DebugDepth:", logText)
	DebugfDepthTrack("track-id-4", 1, "DebugDepthf: %s", logText)
}

func test_Info() {
	Info("Info:", logText)
	Infof("Infof: %s", logText)
	InfoDepth(0, "InfoDepth:", logText)
	InfofDepth(1, "InfoDepthf: %s", logText)
	InfoTrack("track-id-1", "Info:", logText)
	InfofTrack("track-id-2", "Infof: %s", logText)
	InfoDepthTrack("track-id-3", 0, "InfoDepth:", logText)
	InfofDepthTrack("track-id-4", 1, "InfoDepthf: %s", logText)
}

func test_Warn() {
	Warn("Warn:", logText)
	Warnf("Warnf: %s", logText)
	WarnDepth(0, "WarnDepth:", logText)
	WarnfDepth(1, "WarnDepthf: %s", logText)
	WarnTrack("track-id-1", "Warn:", logText)
	WarnfTrack("track-id-2", "Warnf: %s", logText)
	WarnDepthTrack("track-id-3", 0, "WarnDepth:", logText)
	WarnfDepthTrack("track-id-4", 1, "WarnDepthf: %s", logText)
}

func test_Error() {
	Error("Error:", logText)
	Errorf("Errorf: %s", logText)
	ErrorDepth(0, "ErrorDepth:", logText)
	ErrorfDepth(1, "ErrorDepthf: %s", logText)
	ErrorTrack("track-id-1", "Error:", logText)
	ErrorfTrack("track-id-2", "Errorf: %s", logText)
	ErrorDepthTrack("track-id-3", 0, "ErrorDepth:", logText)
	ErrorfDepthTrack("track-id-4", 1, "ErrorDepthf: %s", logText)
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
