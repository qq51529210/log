# 日志库
标准库的log包输出日志头部日期格式固定了（2018/12/29 23:26:15.927321），我不喜欢（我喜欢“2018-12-29”这种格式），所以造了个轮子，可以自由的打印日志头部，包括：级别、时间、堆栈路径。

## 用法

引入包

```go
import "github.com/qq51529210/log"

// 从缓存池中获取*Log，
l := log.GetLog()
// 放回缓存池
defer log.PutLog(l)
// 或者new一个
l = new(log.Log)
// 时间
t := time.Now()
l.Time(&t)
l.Space()
// 级别
l.Level(log.LevelDebug)
l.Space()
// 堆栈
l.FileLine(0)
l.PathLine(0)
l.Space()
// 换行
l.EndLine()
// 上面是自己的格式，或者按我的格式:level time file-line log
l.Header(log.LevelDebug, 0)
// 按我的格式输出日志到控制台，下面输出
// D 2020-12-30 01:17:31.206716000 /Users/ben/Documents/project/go/src/github.com/qq51529210/log/log.go:299 test
l.Print(os.Stdout, log.LevelDebug, 0, "test")
// 格式化输出
l.Printf(os.Stdout, log.LevelDebug, 0, "test %d", 123)
l.Fprint(os.Stdout, log.LevelDebug, 0, 1, 2, 3, 4, 5)
// 也可以用全局函数，但是会输出到默认的io.Writer
log.Print(log.LevelDebug, 0, "test")
log.Printf(log.LevelDebug, 0, "test %d", 123)
log.Fprint(log.LevelDebug, 0, 1, 2, 3, 4, 5)
// 设置全局io.Writer
log.SetWriter(os.Stderr)
// 自定义符号
log.SpaceSeparator    byte      = ' '         // 空格
log.DateSeparator     byte      = '-'         // 日期
log.TimeSeparator     byte      = ':'         // 时间
log.NanoSecSeparator  byte      = '.'         // 纳秒
log.FileLineSeparator byte      = ':'         // 堆栈
// 另外，还实现了一个本地文件的日志File
file, err := log.NewFileLogger(&FileConfig{})
check(err)
defer file.Close()
// 然后，设置Writer
log.SetWriter(file)
```

## 下一步

实现kafka和flume的功能。

## 测试  

下面是和标准库的性能测试，重新造的轮子还是要快一些，然并卵。
```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_LoggerPrint-4          1000000              1018 ns/op             216 B/op          2 allocs/op
Benchmark_StdLoggerPrint-4        982530              1136 ns/op             224 B/op          3 allocs/op
Benchmark_Print-4                1000000              1046 ns/op             216 B/op          2 allocs/op
Benchmark_StdPrint-4             1067096              1131 ns/op             224 B/op          3 allocs/op
Benchmark_Printf-4                995175              1163 ns/op             224 B/op          3 allocs/op
Benchmark_StdPrintf-4            1000000              1202 ns/op             240 B/op          3 allocs/op
Benchmark_Sprint-4               1000000              1139 ns/op             224 B/op          3 allocs/op
Benchmark_StdSprint-4            1000000              1180 ns/op             232 B/op          3 allocs/op
PASS
ok      github.com/qq51529210/log       12.809s
```

