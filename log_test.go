package log

import (
	"bytes"
	"log"
	"testing"
	"time"
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

func TestLog_DateTime(t *testing.T) {
	l.Reset()
	tm := time.Date(20, 1, 12, 2, 30, 20, 123456, time.Local)
	l.Time(&tm)
	if string(l.b) != "0020-01-12 02:30:20.123456000" {
		t.FailNow()
	}
}

func TestLog_Level(t *testing.T) {
	l.Reset()
	l.Level(LevelDebug)
	if l.b[0] != byte(LevelDebug) {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelInfo)
	if l.b[0] != byte(LevelInfo) {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelWarn)
	if l.b[0] != byte(LevelWarn) {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelError)
	if l.b[0] != byte(LevelError) {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelPanic)
	if l.b[0] != byte(LevelPanic) {
		t.FailNow()
	}
	l.Reset()
}

func TestLog_Space(t *testing.T) {
	l.Reset()
	l.Space()
	if string(l.b) != " " {
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
		_, _ = l.Print(&w, LevelDebug, 0, "test")
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
		_, _ = Print(LevelDebug, 0, "test")
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
		_, _ = Printf(LevelDebug, 0, "test%d", i)
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
		_, _ = Fprint(LevelDebug, 0, i)
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
