package log

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

var (
	str = "1234567890"
)

func TestPanic(t *testing.T) {
	defer Recover(os.Stderr, true, true, func() {

	})
	Panic("test panic")
}

func TestPanicStd(t *testing.T) {
	defer Recover(os.Stderr, true, true, func() {

	})
	panic("test std panic")
}

func TestLog(t *testing.T) {
	level := []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelPanic,
	}
	// 打印，级别，堆栈，调用方法就行，并发安全的。
	// 这几个函数是有换行的
	for i := 0; i < 5; i++ {
		Print(os.Stderr, level[i], i, FileLine(i%3), str)
		Printf(os.Stderr, level[i], i, FileLine(i%3), "Printf: %s", str)
		Sprint(os.Stderr, LevelDebug, i, FileLineName, "Sprint: ", str)
	}
	// 不满意我实现的Print，可以自己定义格式，但是换行要自己打印
	logger := Get()
	// 从缓存池里拿出来的，先清空原来的缓存
	logger.Reset()
	// 1.打印级别
	logger.Level(LevelPanic)
	// 2.打印调用堆栈
	logger.FilePathLine(0, FileLineName)
	// 3.打印时间，格式是
	logger.DateTime(6)
	// 4.打印文本
	logger.String(str)
	// 换行
	logger.EndLine()
	// 输出
	os.Stderr.Write(logger.Bytes())
	// 封装的Print方法
	logger.D(os.Stderr, FileLineFullPath, str)
	logger.I(os.Stderr, FileLineFullPath, str)
	logger.W(os.Stderr, FileLineFullPath, str)
	logger.E(os.Stderr, FileLineFullPath, str)
	Put(logger)
	// 设置自己的分隔符
	// 日期格式 2006-01-02
	// 时间格式 15:04:05.999999999
	DateSeparator = '#'
	TimeSeparator = '*'
	NanoSecSeparator = '>'
	SpaceSeparator = '_'
	FileLineSeparator = '|'
	Print(os.Stderr, LevelDebug, 0, FileLineFullPath, str)
}

func Benchmark_StdLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	logger := log.New(buf, "", log.Llongfile|log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Println(str)
		buf.Reset()
	}
}

func Benchmark_MyLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		Print(buf, LevelDebug, 0, FileLineFullPath, str)
		buf.Reset()
	}
}

func Benchmark_Fmt_StdLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	logger := log.New(buf, "d", log.Llongfile|log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		logger.Printf("Printf %s", str)
		buf.Reset()
	}
}

func Benchmark_Fmt_MyLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		Printf(buf, LevelDebug, 0, FileLineFullPath, "Printf %s", str)
		buf.Reset()
	}
}

func Benchmark_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	logger := &Log{}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		logger.Print(buf, LevelDebug, 0, FileLineFullPath, str)
		buf.Reset()
	}
}

func Benchmark_fmt_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	logger := &Log{}
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		logger.Printf(buf, LevelDebug, 0, FileLineFullPath, "Printf %s", str)
		buf.Reset()
	}
}

func testPanicStd(w io.Writer) {
	defer Recover(w, true, true, func() {

	})
	panic("test std panic")
}

func testPanic(w io.Writer) {
	defer Recover(w, true, true, func() {

	})
	Panic("test std panic")
}

func Benchmark_Panic(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		testPanic(buf)
		buf.Reset()
	}
}

func Benchmark_PanicStd(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < b.N; i++ {
		testPanicStd(buf)
		buf.Reset()
	}
}
