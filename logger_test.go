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
	logger.Fprint("logger.Fprint ", 1)
	logger.Fprintf("logger.Fprintf %d", 2)
	logger.PrintStack(0, "logger.PrintStack")
	logger.FprintStack(0, "logger.FprintStack ", 1)
	logger.FprintfStack(0, "logger.FprintfStack %d", 2)
	stdLogger.Println("stdLogger.Println ", 1)
	stdLogger.Printf("stdLogger.Printf %d\n", 2)

	logger.Debug("logger.Debug ", 1)
	logger.DebugStack(0, "logger.DebugDebugStack ", 2)
	logger.Info("logger.Info ", 1)
	logger.InfoStack(0, "logger.InfoStack ", 2)
	logger.Warn("logger.Warn ", 1)
	logger.WarnStack(0, "logger.WarnStack", 2)
	logger.Error("logger.Error ", 1)
	logger.ErrorStack(0, "logger.ErrorStack", 2)

	Print("Print")
	Fprint("Printf ", 2)
	Fprintf("Fprintf %d", 2)
	Debug("Debug ", 1)
	DebugStack(0, "DebugStack ", 2)
	Info("Info ", 1)
	InfoStack(0, "InfoStack ", 2)
	Warn("Warn ", 1)
	WarnStack(0, "WarnStack ", 2)
	Error("Error ", 1)
	ErrorStack(0, "ErrorStack ", 2)
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
