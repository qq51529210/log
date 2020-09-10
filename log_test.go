package log

import (
	"bytes"
	"log"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var l = new(Logger)

func TestLog_Reset(t *testing.T) {
	l.Reset()
	l.WriteString("123")
	l.Reset()
	if len(l.b) > 0 {
		t.FailNow()
	}
}

func TestLog_WriteByte(t *testing.T) {
	l.Reset()
	l.WriteByte('a')
	if l.b[0] != 'a' {
		t.FailNow()
	}
}

func TestLog_WriteEndLine(t *testing.T) {
	l.Reset()
	l.WriteEndLine()
	if l.b[0] != '\n' {
		t.FailNow()
	}
}

func TestLog_WriteInt(t *testing.T) {
	l.Reset()
	l.WriteInt(123)
	if string(l.b) != "123" {
		t.FailNow()
	}

	l.Reset()
	l.WriteInt(-123)
	if string(l.b) != "-123" {
		t.FailNow()
	}
}

func TestLog_WriteIntL0(t *testing.T) {
	l.Reset()
	l.WriteIntL0(123, 7)
	if string(l.b) != "0000123" {
		t.FailNow()
	}
	l.Reset()
	l.WriteIntL0(-123, 7)
	if string(l.b) != "-0000123" {
		t.FailNow()
	}
}

func TestLog_WriteIntR0(t *testing.T) {
	l.Reset()
	l.WriteIntR0(123, 6)
	if string(l.b) != "123000" {
		t.FailNow()
	}
	l.Reset()
	l.WriteIntR0(-123, 6)
	if string(l.b) != "-123000" {
		t.FailNow()
	}
}

func TestLog_WriteDateTime(t *testing.T) {
	l.Reset()
	tm := time.Date(20, 1, 12, 2, 30, 20, 123456789, time.Local)
	l.WriteDateTime(&tm)
	if string(l.b) != "0020-01-12 02:30:20.123456" {
		t.FailNow()
	}
}

func TestLog_WriteLevel(t *testing.T) {
	l.Reset()
	l.WriteLevel(LevelDebug)
	if l.b[0] != byte(LevelDebug) {
		t.FailNow()
	}
	l.Reset()
	l.WriteLevel(LevelInfo)
	if l.b[0] != byte(LevelInfo) {
		t.FailNow()
	}
	l.Reset()
	l.WriteLevel(LevelWarn)
	if l.b[0] != byte(LevelWarn) {
		t.FailNow()
	}
	l.Reset()
	l.WriteLevel(LevelError)
	if l.b[0] != byte(LevelError) {
		t.FailNow()
	}
	l.Reset()
	l.WriteLevel(LevelPanic)
	if l.b[0] != byte(LevelPanic) {
		t.FailNow()
	}
	l.Reset()
	l.WriteLevel(LevelRecover)
	if l.b[0] != byte(LevelRecover) {
		t.FailNow()
	}
}

func TestLog_WriteStack(t *testing.T) {
	f := "/project/main.go"
	ln := 123
	l.Reset()
	l.WriteStackFile(f, ln)
	_, ff := filepath.Split(f)
	if string(l.b) != ff+":"+strconv.Itoa(ln) {
		t.FailNow()
	}
	l.Reset()
	l.WriteStackPath(f, ln)
	if string(l.b) != f+":"+strconv.Itoa(ln) {
		t.FailNow()
	}
}

func TestLog_WriteSpace(t *testing.T) {
	l.Reset()
	l.WriteSpace()
	if string(l.b) != " " {
		t.FailNow()
	}
}

func TestLog_WriteString(t *testing.T) {
	l.Reset()
	l.WriteString("123")
	if string(l.b) != "123" {
		t.FailNow()
	}
}

func TestLog_Stack(t *testing.T) {
	l.Reset()
	defer func() {
		recover()
		l.WriteStack(false)
		t.Log(string(l.b))
	}()
	panic("123")
}

func Benchmark_LoggerPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		l.Reset()
		_, _ = l.Print(&w, LevelDebug, StackInfoDisable, 0, "test")
		w.Reset()
	}
}

func Benchmark_StdLoggerPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	l := log.New(&w, "D", log.LstdFlags|log.Lmicroseconds)
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
	SetStackInfo(StackInfoDisable)
	for i := 0; i < b.N; i++ {
		_, _ = Print(LevelDebug, 0, "test")
		w.Reset()
	}
}

func Benchmark_StdPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
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
	SetStackInfo(StackInfoDisable)
	for i := 0; i < b.N; i++ {
		_, _ = Printf(LevelDebug, 0, "test%d", i)
		w.Reset()
	}
}

func Benchmark_StdPrintf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
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
	SetStackInfo(StackInfoDisable)
	for i := 0; i < b.N; i++ {
		_, _ = Sprint(LevelDebug, 0, i)
		w.Reset()
	}
}

func Benchmark_StdSprint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(&w)
	for i := 0; i < b.N; i++ {
		log.Println(i)
		w.Reset()
	}
}
