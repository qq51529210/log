package log

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var (
	benchmarkString = "1234567890asdfghjklqwertyuiopmnbvcxz,./';][=-<>?:{}|\\+_)(*&^%$#@!~"
)

func TestLog(t *testing.T) {
	level := []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelPanic,
	}
	for i := 0; i < 5; i++ {
		Print(os.Stderr, level[i], 0, FileLine(i%3), benchmarkString)
		Printf(os.Stderr, level[i], 0, FileLine(i%3), "Printf: %s", benchmarkString)
	}
	logger := &Log{}
	logger.D(os.Stderr, FileLineFullPath, benchmarkString)
	logger.I(os.Stderr, FileLineFullPath, benchmarkString)
	logger.W(os.Stderr, FileLineFullPath, benchmarkString)
	logger.E(os.Stderr, FileLineFullPath, benchmarkString)
	logger.Sprint(os.Stderr, LevelDebug, 0, FileLineName, "Sprint: ", benchmarkString)
}

func Benchmark_StdLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	logger := log.New(buf, "", log.Llongfile|log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Println(benchmarkString)
		buf.Reset()
	}
}

func Benchmark_MyLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		Print(buf, LevelDebug, 0, FileLineFullPath, benchmarkString)
		buf.Reset()
	}
}

func Benchmark_Fmt_StdLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	logger := log.New(buf, "d", log.Llongfile|log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Printf("Printf %s", benchmarkString)
		buf.Reset()
	}
}

func Benchmark_Fmt_MyLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		Printf(buf, LevelDebug, 0, FileLineFullPath, "Printf %s", benchmarkString)
		buf.Reset()
	}
}

func Benchmark_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	logger := &Log{}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		logger.Print(buf, LevelDebug, 0, FileLineFullPath, benchmarkString)
	}
}

func Benchmark_fmt_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	logger := &Log{}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		logger.Printf(buf, LevelDebug, 0, FileLineFullPath, "Printf %s", benchmarkString)
	}
}

func Benchmark_MyLogD(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	logger := &Log{}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		logger.D(buf, FileLineFullPath, benchmarkString)
	}
}
