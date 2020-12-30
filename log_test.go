package log

import (
	"bytes"
	"log"
	"runtime"
	"testing"
)

var (
	l = new(Log)
	w = bytes.Buffer{}
	//stdFlag = 0
	stdFlag = log.LstdFlags | log.Lmicroseconds | log.Llongfile
	//stdFlag = log.LstdFlags | log.Lmicroseconds
)

func TestLog_Reset(t *testing.T) {
	l.Reset()
	l.String("123")
	l.Reset()
	if len(l.b) > 0 {
		t.FailNow()
	}
}

func TestLog_Byte(t *testing.T) {
	l.Reset()
	l.Byte('a')
	if l.b[0] != 'a' {
		t.FailNow()
	}
}

func TestLog_EndLine(t *testing.T) {
	l.Reset()
	l.EndLine()
	if l.b[0] != '\n' {
		t.FailNow()
	}
}

func TestLog_Int(t *testing.T) {
	l.Reset()
	l.Int(123)
	if string(l.b) != "123" {
		t.FailNow()
	}

	l.Reset()
	l.Int(-123)
	if string(l.b) != "-123" {
		t.FailNow()
	}
}

func TestLog_IntL0(t *testing.T) {
	l.Reset()
	l.IntL0(123, 7)
	if string(l.b) != "0000123" {
		t.FailNow()
	}
	l.Reset()
	l.IntL0(-123, 7)
	if string(l.b) != "-0000123" {
		t.FailNow()
	}
}

func TestLog_IntR0(t *testing.T) {
	l.Reset()
	l.IntR0(123, 6)
	if string(l.b) != "123000" {
		t.FailNow()
	}
	l.Reset()
	l.IntR0(-123, 6)
	if string(l.b) != "-123000" {
		t.FailNow()
	}
}

func TestLog_String(t *testing.T) {
	l.Reset()
	l.String("123")
	if string(l.b) != "123" {
		t.FailNow()
	}
}

func Benchmark_LoggerPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		_, path, line, o := runtime.Caller(0)
		if !o {
			path = "???"
			line = -1
		}
		l.Reset()
		l.Header(DebugLevel, path, line)
		l.String("test\n")
		w.Write(l.b)
	}
}

func Benchmark_StdLoggerPrint(b *testing.B) {
	l := log.New(&w, "D", stdFlag)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		l.Output(0, "test")
	}
}

func Benchmark_Print(b *testing.B) {
	SetWriter(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		Print(DebugLevel, BaseSkip, "test")
	}
}

func Benchmark_StdPrint(b *testing.B) {
	log.SetFlags(stdFlag)
	log.SetOutput(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		log.Println("test")
	}
}

func Benchmark_Printf(b *testing.B) {
	SetWriter(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		Printf(DebugLevel, 1, "test%d", i)
	}
}

func Benchmark_StdPrintf(b *testing.B) {
	log.SetFlags(stdFlag)
	log.SetOutput(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		log.Printf("test%d\n", i)
	}
}

func Benchmark_Fprint(b *testing.B) {
	SetWriter(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		Fprint(DebugLevel, BaseSkip, i)
	}
}

func Benchmark_StdSprint(b *testing.B) {
	log.SetFlags(stdFlag)
	log.SetOutput(&w)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Reset()
		log.Println(i)
	}
}
