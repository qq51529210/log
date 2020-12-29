package log

import (
	"bytes"
	"log"
	"testing"
)

var l = new(Log)

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
	w := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		l.Reset()
		l.Header(DebugLevel, 0)
		l.String("test\n")
		w.Reset()
	}
}

func Benchmark_StdLoggerPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	l := log.New(&w, "D", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	for i := 0; i < b.N; i++ {
		l.Println("test")
		w.Reset()
	}
}

func Benchmark_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	SetWriter(&w)
	for i := 0; i < b.N; i++ {
		Print(DebugLevel, 0, "test")
		w.Reset()
	}
}

func Benchmark_StdPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(&w)
	for i := 0; i < b.N; i++ {
		log.Println("test")
		w.Reset()
	}
}

func Benchmark_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	SetWriter(&w)
	for i := 0; i < b.N; i++ {
		Printf(DebugLevel, 0, "test%d", i)
		w.Reset()
	}
}

func Benchmark_StdPrintf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(&w)
	for i := 0; i < b.N; i++ {
		log.Printf("test%d\n", i)
		w.Reset()
	}
}

func Benchmark_Sprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	SetWriter(&w)
	for i := 0; i < b.N; i++ {
		Fprint(DebugLevel, 0, i)
		w.Reset()
	}
}

func Benchmark_StdSprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(&w)
	for i := 0; i < b.N; i++ {
		log.Println(i)
		w.Reset()
	}
}
