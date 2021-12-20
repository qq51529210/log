# log

Standard library log format header is fixed(2006/01/02 15:04:05.123456),I don't like it,so I wrote this package.

You can output your own header.

## Usage

```go
Debug("Debug:", logText)
Debugf("Debugf: %s", logText)
DepthDebug(1, "DepthDebug:", logText)
DepthDebugf(1, "DepthDebugf: %s", logText)
LevelDepthDebug(0, 0, "LevelDepthDebug:", logText)
LevelDepthDebug(-1, 1, "LevelDepthDebug:", logText)
LevelDepthDebugf(0, 0, "LevelDepthDebugf: %s", logText)
LevelDepthDebugf(-1, 1, "LevelDepthDebugf: %s", logText)
// 
myLogger := NewLogger(os.Stderr, 1, new(CallStackFilePathHeader))
myLogger.SetWriter(myStrogeWriter)
myLogger.SetLevel(myLevel)
myLogger.SetType(OutputTypeDebug, OutputTypeInfo, OutputTypeWarn, OutputTypeError)
myLogger.Debug("Debug:", logText)
myLogger.Debugf("Debugf: %s", logText)
myLogger.DepthDebug(1, "DepthDebug:", logText)
myLogger.DepthDebugf(1, "DepthDebugf: %s", logText)
myLogger.myLogger.LevelDepthDebug(0, 0, "LevelDepthDebug:", logText)
myLogger.LevelDepthDebug(-1, 1, "LevelDepthDebug:", logText)
myLogger.LevelDepthDebugf(0, 0, "LevelDepthDebugf: %s", logText)
myLogger.LevelDepthDebugf(-1, 1, "LevelDepthDebugf: %s", logText)
// 
writer := NewFile(logDir, 1024*1024, 7, time.Second)
SetWriter(writer)
SetLevel(1))
// will not output
LevelDebug(0, "LevelDebug")
// output LevelDebug
LevelDebug(1, "LevelDebug")
// output LevelDebug
LevelDebug(2, "LevelDebug")
SetType(OutputTypeInfo, OutputTypeError)
// will not output
Debug("Debug")
// output Info
Info("Info")
// will not output
Warn("Warn")
// output Error
Error("Error")
```

## Writer

- [File](./file.go)

  Save log on local disk.

## Benchmark  

Comparison with standard library.

```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_My_Logger-4             763563              1375 ns/op             248 B/op          4 allocs/op
Benchmark_Std_Logger-4            991500              1212 ns/op             240 B/op          4 allocs/op
Benchmark_My_Logger_f-4           860084              1399 ns/op             248 B/op          4 allocs/op
Benchmark_Std_Logger_f-4          955314              1262 ns/op             248 B/op          4 allocs/op
PASS
ok      github.com/qq51529210/log       8.100s
```

