package log

import (
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	lg := NewLogger(os.Stdout, new(DefaultHeader), "default")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
	//
	lg = NewLogger(os.Stdout, new(FileNameHeader), "file name")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
	//
	lg = NewLogger(os.Stdout, new(FilePathHeader), "file path")
	testLoggerDebug(t, lg)
	testLoggerInfo(t, lg)
	testLoggerWarn(t, lg)
	testLoggerError(t, lg)
}

func testLoggerDebug(t *testing.T, lg Logger) {
	lg.Debug(1, 2, 3)
	lg.Debugf("%d", 1)
	lg.DebugDepth(0, 1, 2, 3)
	lg.DebugfDepth(0, "%d", 1)
	lg.DebugTrace("trace", 1, 2)
	lg.DebugfTrace("trace", "%d", 1)
	lg.DebugDepthTrace(0, "trace", 1, 2)
	lg.DebugfDepthTrace(0, "trace", "%d", 1)
}

func testLoggerInfo(t *testing.T, lg Logger) {
	lg.Info(1, 2, 3)
	lg.Infof("%d", 1)
	lg.InfoDepth(0, 1, 2, 3)
	lg.InfofDepth(0, "%d", 1)
	lg.InfoTrace("trace", 1, 2)
	lg.InfofTrace("trace", "%d", 1)
	lg.InfoDepthTrace(0, "trace", 1, 2)
	lg.InfofDepthTrace(0, "trace", "%d", 1)
}

func testLoggerWarn(t *testing.T, lg Logger) {
	lg.Warn(1, 2, 3)
	lg.Warnf("%d", 1)
	lg.WarnDepth(0, 1, 2, 3)
	lg.WarnfDepth(0, "%d", 1)
	lg.WarnTrace("trace", 1, 2)
	lg.WarnfTrace("trace", "%d", 1)
	lg.WarnDepthTrace(0, "trace", 1, 2)
	lg.WarnfDepthTrace(0, "trace", "%d", 1)
}

func testLoggerError(t *testing.T, lg Logger) {
	lg.Error(1, 2, 3)
	lg.Errorf("%d", 1)
	lg.ErrorDepth(0, 1, 2, 3)
	lg.ErrorfDepth(0, "%d", 1)
	lg.ErrorTrace("trace", 1, 2)
	lg.ErrorfTrace("trace", "%d", 1)
	lg.ErrorDepthTrace(0, "trace", 1, 2)
	lg.ErrorfDepthTrace(0, "trace", "%d", 1)
}

func Test_Recover(t *testing.T) {
	lg := NewLogger(os.Stdout, new(DefaultHeader), "default")
	defer func() {
		lg.Recover(recover())
	}()

	testRecover()
}

func testRecover() {
	testRecover1()
}

func testRecover1() {
	testRecover2()
}

func testRecover2() {
	panic("test recover")
}

// func Benchmark_My_Logger(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	var output bytes.Buffer
// 	myLogger.SetOutput(&output)
// 	for i := 0; i < b.N; i++ {
// 		myLogger.Debug(logText)
// 		output.Reset()
// 	}
// }

// func Benchmark_Std_Logger(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	var output bytes.Buffer
// 	stdLogger.SetOutput(&output)
// 	for i := 0; i < b.N; i++ {
// 		stdLogger.Println(logText)
// 		output.Reset()
// 	}
// }

// func Benchmark_My_Logger_f(b *testing.B) {
// 	myLogger.Info("benchmark my logger")
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	var output bytes.Buffer
// 	myLogger.SetOutput(&output)
// 	for i := 0; i < b.N; i++ {
// 		myLogger.Debugf(logFormat, logText)
// 		output.Reset()
// 	}
// }

// func Benchmark_Std_Logger_f(b *testing.B) {
// 	stdLogger.Println("benchmark std logger")
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	var output bytes.Buffer
// 	stdLogger.SetOutput(&output)
// 	for i := 0; i < b.N; i++ {
// 		stdLogger.Printf(logFormat, logText)
// 		output.Reset()
// 	}
// }
