package log

import "testing"

// import (
// 	"bytes"
// 	"log"
// 	"runtime"
// 	"testing"
// )

// var (
// 	l = new(Log)
// 	w = bytes.Buffer{}
// 	//stdFlag = 0
// 	stdFlag = log.LstdFlags | log.Lmicroseconds | log.Llongfile
// 	//stdFlag = log.LstdFlags | log.Lmicroseconds
// )

func TestLog_WriteRightAlignInt(t *testing.T) {
	log := GetLog()
	defer PutLog(log)

	log.WriteRightAlignInt(1234, 7)
	if log.String() != "0001234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(-1234, 3)
	if log.String() != "-1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(-1234, 6)
	if log.String() != "-001234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(0, 4)
	if log.String() != "0000" {
		t.FailNow()
	}
}

func TestLog_WriteLeftAlignInt(t *testing.T) {
	log := GetLog()
	defer PutLog(log)

	log.WriteLeftAlignInt(1234, 7)
	if log.String() != "1234000" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(-1234, 3)
	if log.String() != "-1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(-1234, 6)
	if log.String() != "-123400" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(0, 4)
	if log.String() != "0000" {
		t.FailNow()
	}
}

func TestLog_WriteInt(t *testing.T) {
	log := GetLog()
	defer PutLog(log)

	log.WriteInt(1234)
	if log.String() != "1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteInt(-1234)
	if log.String() != "-1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteInt(0)
	if log.String() != "0" {
		t.FailNow()
	}
}

// func TestLog_Byte(t *testing.T) {
// 	l.Reset()
// 	l.Byte('a')
// 	if l.b[0] != 'a' {
// 		t.FailNow()
// 	}
// }

// func TestLog_EndLine(t *testing.T) {
// 	l.Reset()
// 	l.EndLine()
// 	if l.b[0] != '\n' {
// 		t.FailNow()
// 	}
// }

// func TestLog_Int(t *testing.T) {
// 	l.Reset()
// 	l.Int(123)
// 	if string(l.b) != "123" {
// 		t.FailNow()
// 	}

// 	l.Reset()
// 	l.Int(-123)
// 	if string(l.b) != "-123" {
// 		t.FailNow()
// 	}
// }

// func TestLog_IntL0(t *testing.T) {
// 	l.Reset()
// 	l.IntL0(123, 7)
// 	if string(l.b) != "0000123" {
// 		t.FailNow()
// 	}
// 	l.Reset()
// 	l.IntL0(-123, 7)
// 	if string(l.b) != "-0000123" {
// 		t.FailNow()
// 	}
// }

// func TestLog_IntR0(t *testing.T) {
// 	l.Reset()
// 	l.IntR0(123, 6)
// 	if string(l.b) != "123000" {
// 		t.FailNow()
// 	}
// 	l.Reset()
// 	l.IntR0(-123, 6)
// 	if string(l.b) != "-123000" {
// 		t.FailNow()
// 	}
// }

// func TestLog_String(t *testing.T) {
// 	l.Reset()
// 	l.String("123")
// 	if string(l.b) != "123" {
// 		t.FailNow()
// 	}
// }

// func TestLog(t *testing.T) {
// 	Print(DebugLevel, 0, "test")
// 	Debug("test")
// }

// func Benchmark_LoggerPrint(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		_, path, line, o := runtime.Caller(0)
// 		if !o {
// 			path = "???"
// 			line = -1
// 		}
// 		l.Reset()
// 		defaultPrintHeader(l, DebugLevel, path, line)
// 		l.String("test\n")
// 		w.Write(l.b)
// 	}
// }

// func Benchmark_StdLoggerPrint(b *testing.B) {
// 	l := log.New(&w, "D", stdFlag)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		l.Output(0, "test")
// 	}
// }

// func Benchmark_Print(b *testing.B) {
// 	SetWriter(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		Print(DebugLevel, 0, "test")
// 	}
// }

// func Benchmark_StdPrint(b *testing.B) {
// 	log.SetFlags(stdFlag)
// 	log.SetOutput(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		log.Println("test")
// 	}
// }

// func Benchmark_Printf(b *testing.B) {
// 	SetWriter(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		Printf(DebugLevel, 1, "test%d", i)
// 	}
// }

// func Benchmark_StdPrintf(b *testing.B) {
// 	log.SetFlags(stdFlag)
// 	log.SetOutput(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		log.Printf("test%d\n", i)
// 	}
// }

// func Benchmark_Fprint(b *testing.B) {
// 	SetWriter(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		Fprint(DebugLevel, 0, i)
// 	}
// }

// func Benchmark_StdSprint(b *testing.B) {
// 	log.SetFlags(stdFlag)
// 	log.SetOutput(&w)
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		log.Println(i)
// 	}
// }
