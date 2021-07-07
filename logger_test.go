package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
	logger    = NewLogger(os.Stderr, nil, nil)
)

func init() {
	logger.fmtTimeHeader = FormatTimeHeader
	logger.fmtStackHeader = FormatFilePathStackHeader
	// logger.callerHeaderFunc = PrintFileNameCallerHeader
	stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	// stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func Test_Logger(t *testing.T) {
	// Check out std output.
	logger.Print("logger.Print")
	logger.Fprint("logger.Fprint", 1, 1)
	logger.Fprintf("logger.Fprintf %d", 2)
	//
	logger.PrintStack(0, "logger.PrintStack")
	logger.FprintStack(0, "logger.FprintStack", 1, 1)
	logger.FprintfStack(0, "logger.FprintfStack %d", 2)
	//
	stdLogger.Println("stdLogger.Println", 1)
	stdLogger.Printf("stdLogger.Printf %d\n", 2)
	//
	logger.Debug("logger.Debug", 1, 1)
	logger.DebugStack(0, "logger.DebugStack", 1, 2)
	logger.Debugf("logger.Debugf %d", 1)
	logger.DebugfStack(0, "logger.DebugfStack %d", 2)
	//
	logger.Info("logger.Info", 1, 1)
	logger.InfoStack(0, "logger.InfoStack", 1, 2)
	logger.Infof("logger.Infof %d", 1)
	logger.InfofStack(0, "logger.InfofStack %d", 2)
	//
	logger.Warn("logger.Warn", 1, 1)
	logger.WarnStack(0, "logger.WarnStack", 1, 2)
	logger.Warnf("logger.Warnf %d", 1)
	logger.WarnfStack(0, "logger.WarnfStack %d", 2)
	//
	logger.Error("logger.Error", 1, 1)
	logger.ErrorStack(0, "logger.ErrorStack", 1, 2)
	logger.Errorf("logger.Errorf %d", 1)
	logger.ErrorfStack(0, "logger.ErrorfStack %d", 2)
	//
	Print("Print")
	Fprint("Printf", 1, 2)
	Fprintf("Fprintf %d", 2)
	//
	Debug("Debug", 1, 1)
	DebugStack(0, "DebugStack", 1, 2)
	Debugf("Debugf %d", 1)
	DebugfStack(0, "DebugfStack %d", 2)
	//
	Info("Info", 1, 1)
	InfoStack(0, "InfoStack", 1, 2)
	Infof("Infof %d", 1)
	InfofStack(0, "InfofStack %d", 2)
	//
	Warn("Warn", 1, 1)
	WarnStack(0, "WarnStack", 1, 2)
	Warnf("Warnf %d", 1)
	WarnfStack(0, "WarnfStack %d", 2)
	//
	Error("Error ", 1, 1)
	ErrorStack(0, "ErrorStack", 1, 2)
	Errorf("Errorf %d", 1)
	ErrorfStack(0, "ErrorfStack %d", 2)
}

func Benchmark_Logger_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.out = &output
	for i := 0; i < b.N; i++ {
		logger.Print("Print")
		output.Reset()
	}
}

func Benchmark_Logger_Fprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.out = &output
	for i := 0; i < b.N; i++ {
		logger.Fprint("Print ", 1)
		output.Reset()
	}
}

func Benchmark_StdLogger_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Println("Print ", 1)
		output.Reset()
	}
}

func Benchmark_Logger_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.out = &output
	for i := 0; i < b.N; i++ {
		logger.Fprintf("Printf %d", 2)
		output.Reset()
	}
}

func Benchmark_StdLogger_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Printf("Printf %d\n", 2)
		output.Reset()
	}
}
