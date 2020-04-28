package log

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

func TestLog_Reset(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.String("123")
	l.Reset()
	if len(l.b) > 0 {
		t.FailNow()
	}
}

func TestLog_IntegerAlignLeft(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.IntegerAlignLeft(123, 7)
	if string(l.b) != "0000123" {
		t.FailNow()
	}
}

func TestLog_IntegerAlignRight(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.IntegerAlignRight(123, 6)
	if string(l.b) != "123000" {
		t.FailNow()
	}
}

func TestLog_Integer(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.Integer(123)
	if string(l.b) != "123" {
		t.FailNow()
	}
}

func TestLog_Byte(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.Byte('a')
	if l.b[0] != 'a' {
		t.FailNow()
	}
}

func TestLog_EndLine(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.EndLine()
	if l.b[0] != '\n' {
		t.FailNow()
	}
}

func TestLog_DateTime(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.DateTime(9)
	t.Log(string(l.b))
}

func TestLog_Level(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.Level(LevelDebug)
	if l.b[0] != 'D' {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelInfo)
	if l.b[0] != 'I' {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelWarn)
	if l.b[0] != 'W' {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelError)
	if l.b[0] != 'E' {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelPanic)
	if l.b[0] != 'P' {
		t.FailNow()
	}
	l.Reset()
	l.Level(LevelRecover)
	if l.b[0] != 'R' {
		t.FailNow()
	}
}

func TestLog_FilePathLine(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.FilePathLine(1, FileLineDisable)
	t.Log(string(l.b))
	l.Reset()
	l.FilePathLine(1, FileLineFullPath)
	t.Log(string(l.b))
	l.Reset()
	l.FilePathLine(1, FileLineName)
	t.Log(string(l.b))
}

func TestLog_String(t *testing.T) {
	l := Get()
	l.Reset()
	defer Put(l)
	l.String("123")
	if string(l.b) != "123" {
		t.FailNow()
	}
}

func TestLog_Stack(t *testing.T) {
	l := Get()
	l.Reset()
	defer func() {
		recover()
		l.Stack(false)
		t.Log(string(l.b))
		defer Put(l)
	}()
	panic("123")
}

func Test_Panic(t *testing.T) {
	defer func() {
		re := recover()
		switch re.(type) {
		case *panicError:
		default:
			t.FailNow()
		}
	}()
	Panic("123")
}

func Test_CheckError(t *testing.T) {
	defer func() {
		re := recover()
		switch re.(type) {
		case *panicError:
		default:
			t.FailNow()
		}
	}()
	CheckError(errors.New("123"))
}

func Benchmark_LogPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	l := Get()
	defer Put(l)
	l.Info.Writer = &w
	l.Info.Level = LevelDebug
	l.Info.Skip = 0
	l.Info.FileLine = FileLineDisable
	for i := 0; i < b.N; i++ {
		l.Reset()
		l.Print("test")
		w.Reset()
	}
}

func Benchmark_Std_LogPrint(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	l := log.New(&w, "D", log.LstdFlags|log.Lmicroseconds)
	for i := 0; i < b.N; i++ {
		l.Println("test")
		w.Reset()
	}
}

func Benchmark_Std_Print(b *testing.B) {
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

func Benchmark_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		Print(&w, LevelDebug, 0, FileLineDisable, "test")
		w.Reset()
	}
}

func Benchmark_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	w := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		Printf(&w, LevelDebug, 0, FileLineDisable, "test%d", i)
		w.Reset()
	}
}

func Benchmark_Std_Printf(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		Sprint(&w, LevelDebug, 0, FileLineDisable, i)
		w.Reset()
	}
}

func Benchmark_Std_Sprint(b *testing.B) {
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
