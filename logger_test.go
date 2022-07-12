package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	myLogger  = NewLogger(os.Stderr, NewHeaderFormater(FilePathStackHeaderFormater, "appID"))
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
	DebugTrace("trace-id-1", "Debug:", logText)
	DebugfTrace("trace-id-2", "Debugf: %s", logText)
	DebugDepthTrace("trace-id-3", 0, "DebugDepth:", logText)
	DebugfDepthTrace("trace-id-4", 1, "DebugDepthf: %s", logText)
}

func test_Info() {
	Info("Info:", logText)
	Infof("Infof: %s", logText)
	InfoDepth(0, "InfoDepth:", logText)
	InfofDepth(1, "InfoDepthf: %s", logText)
	InfoTrace("trace-id-1", "Info:", logText)
	InfofTrace("trace-id-2", "Infof: %s", logText)
	InfoDepthTrace("trace-id-3", 0, "InfoDepth:", logText)
	InfofDepthTrace("trace-id-4", 1, "InfoDepthf: %s", logText)
}

func test_Warn() {
	Warn("Warn:", logText)
	Warnf("Warnf: %s", logText)
	WarnDepth(0, "WarnDepth:", logText)
	WarnfDepth(1, "WarnDepthf: %s", logText)
	WarnTrace("trace-id-1", "Warn:", logText)
	WarnfTrace("trace-id-2", "Warnf: %s", logText)
	WarnDepthTrace("trace-id-3", 0, "WarnDepth:", logText)
	WarnfDepthTrace("trace-id-4", 1, "WarnDepthf: %s", logText)
}

func test_Error() {
	Error("Error:", logText)
	Errorf("Errorf: %s", logText)
	ErrorDepth(0, "ErrorDepth:", logText)
	ErrorfDepth(1, "ErrorDepthf: %s", logText)
	ErrorTrace("trace-id-1", "Error:", logText)
	ErrorfTrace("trace-id-2", "Errorf: %s", logText)
	ErrorDepthTrace("trace-id-3", 0, "ErrorDepth:", logText)
	ErrorfDepthTrace("trace-id-4", 1, "ErrorDepthf: %s", logText)
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
