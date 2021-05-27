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
	// logger.PrintCallerHeader = PrintFilePathCallerHeader
	// logger.PrintCallerHeader = PrintFileNameCallerHeader
	// stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	// stdLogger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func Test_Logger(t *testing.T) {
	// Check out std output.
	logger.Print("test ", 1)
	logger.Printf("test %d", 2)
	stdLogger.Println("test ", 1)
	stdLogger.Printf("test %d\n", 2)
}

func Benchmark_Logger_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.Writer = &output
	for i := 0; i < b.N; i++ {
		logger.Print("test ", 1)
		output.Reset()
	}
}

func Benchmark_StdLogger_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Println("test ", 1)
		output.Reset()
	}
}

func Benchmark_Logger_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	logger.Writer = &output
	for i := 0; i < b.N; i++ {
		logger.Printf("test %d", 2)
		output.Reset()
	}
}

func Benchmark_StdLogger_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var output bytes.Buffer
	stdLogger.SetOutput(&output)
	for i := 0; i < b.N; i++ {
		stdLogger.Printf("test %d\n", 2)
		output.Reset()
	}
}
