# log

Standard library log format header is fixed(2006/01/02 15:04:05.123456),I don't like it,so I wrote this package.

You can output your own header.

## How to use

```go
logger := NewLogger(os.StdErr)
// Default format is like "2006-01-02 15:04:05.123456 test"
logger.Print("test")
logger.Printf("test %d", 1)
// Set your own print header function.
logger.PrintTimeHeader = func(*Log){}
// Set your own print caller file path function.
logger.PrintCallerHeader = func(*Log,int){}
// Set PrintFilePathCallerHeader or PrintFileNameCallerHeader
logger.PrintCallerHeader = PrintFilePathCallerHeader
```

## Write

- [File](./file_writer.go)

  Save log on local disk.

- Flume

  Save log to Flume.To be implemented.

- Kafka

  Save log to Kafka.To be implemented.

## Benchmark  

Comparison with standard library.

```go
2021-05-27 18:51:00.409254 test 1
2021-05-27 18:51:00.409372 test 2
2021/05/27 18:51:00.409376 test  1
2021/05/27 18:51:00.409380 test 2
```
```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_Logger_Print-4         2875112               409 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Print-4      2890418               413 ns/op               8 B/op          1 allocs/op
Benchmark_Logger_Printf-4        2980972               394 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Printf-4     2902381               433 ns/op               8 B/op          1 allocs/op
PASS
ok      github.com/qq51529210/log       7.224s
```
```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_Logger_Print-4         2853477               410 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Print-4      2879283               420 ns/op               8 B/op          1 allocs/op
Benchmark_Logger_Printf-4        3014130               395 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Printf-4     2902104               410 ns/op               8 B/op          1 allocs/op
PASS
ok      github.com/qq51529210/log       6.545s
```
```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_Logger_Print-4         2886384               411 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Print-4      2880801               417 ns/op               8 B/op          1 allocs/op
Benchmark_Logger_Printf-4        2973010               398 ns/op               0 B/op          0 allocs/op
Benchmark_StdLogger_Printf-4     2901409               408 ns/op               8 B/op          1 allocs/op
PASS
ok      github.com/qq51529210/log       9.507s
```

