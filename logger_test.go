package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	stdLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
	logger    = NewLogger(os.Stderr)
)

func init() {
	logger.PrintCallerHeader = PrintFilePathCallerHeader
	// logger.PrintCallerHeader = PrintFileNameCallerHeader
	stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	// stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	SetPrintCallerHeader(PrintFilePathCallerHeader)
}

func Test_Logger(t *testing.T) {
	// Check out std output.
	logger.Print("logger.Print ", 1)
	logger.Printf("logger.Printf %d", 2)
	stdLogger.Println("stdLogger.Println ", 1)
	stdLogger.Printf("stdLogger.Printf %d\n", 2)

	logger.Debug("logger.Debug ", 1)
	logger.Debugf("logger.Debugf %d", 2)
	logger.Info("logger.Info ", 1)
	logger.Infof("logger.Infof %d", 2)
	logger.Warn("logger.Warn ", 1)
	logger.Warnf("logger.Warnf %d", 2)
	logger.Error("logger.Error ", 1)
	logger.Errorf("logger.Errorf %d", 2)

	Print("Print ", 1)
	Printf("Printf %d", 2)
	Debug("Debug ", 1)
	Debugf("Debugf %d", 2)
	Info("Info ", 1)
	Infof("Infof %d", 2)
	Warn("Warn ", 1)
	Warnf("Warnf %d", 2)
	Error("Error ", 1)
	Errorf("Errorf %d", 2)
}

func Benchmark_Logger_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.Writer = &output
	for i := 0; i < b.N; i++ {
		logger.Print("Print ", 1)
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
	logger.Writer = &output
	for i := 0; i < b.N; i++ {
		logger.Printf("Printf %d", 2)
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
