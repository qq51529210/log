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
	for i := 0; i < 3; i++ {
		Print(os.Stderr, LevelDebug, 0, FileLine(i), benchmarkString)
		Printf(os.Stderr, LevelDebug, 0, FileLine(i), "output: %s", benchmarkString)
	}
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
		logger.Printf("log %s", benchmarkString)
		buf.Reset()
	}
}

func Benchmark_Fmt_MyLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		Printf(buf, LevelDebug, 0, FileLineFullPath, "log %s", benchmarkString)
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
		logger.Printf(buf, LevelDebug, 0, FileLineFullPath, "log %s", benchmarkString)
	}
}
