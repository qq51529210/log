# 日志库
标准库的log包输出日志头部日期格式固定了（2018/12/29 23:26:15.927321），我不喜欢（我喜欢“2018-12-29”这种格式），所以造了个轮子，可以自由的打印日志头部，包括：级别、时间、堆栈路径。

## 用法

自定义头部格式，比如，`20201230 024308884189 [Debug] /Users/ben/Documents/project/go/src/test/main.go:26 test log string`	

```go
// 首先，设置日期和时间的分隔符为""
log.DateSeparator = ""
log.TimeSeparator = ""
func Print(w io.Writer, level string, str string) (int, error) {
  // 从缓存池中获取*Log，
  l := log.GetLog()
  // 时间
  l.Time()
  l.Space()
  // 级别
  l.String(level)
  l.Space()
  // 调用堆栈
  l.PathLine(1)
  l.Space()
  // 文本
  l.String(str)
  l.EndLine()
  // 输出
  n, err := w.Write(l.Data())
  // 放回缓存池
  log.PutLog(l)
  // 返回
  return n, err
}
```

默认格式，比如，`[D] 2020-12-30 02:53:20.953755 /Users/ben/Documents/project/go/src/test/main.go:20 test`

```go
// 设置默认的io.Writer
log.SetWriter(os.Stderr)
// 设置级别
log.DebugLevel = "[D]"
// 输出
log.Debug("test")
// 如果需要控制调用堆栈
log.Print("debug", 0, "test")
log.Printf("info", 1, "test %d", 0)
log.Fprint("warn", 2, 0, "1", 2.3)
```

保存到本地磁盘

```go
// 实例，FileConfig的字段，请看注释
file, err := NewFileLogger(&FileConfig{})
if err != nil {
  panic(err)
}
// 关闭
defer func(){
  _ = file.Close()
}
// 设置io.Writer
log.SetWriter(file)
// 写的日志会输出到file
log.Debug("test")
```

## 下一步

实现kafka和flume的功能。

## 测试  

下面是和标准库的性能测试，重新造的轮子还是要快一些，然并卵。
```go
goos: darwin
goarch: amd64
Benchmark_LoggerPrint-4          1428115               830 ns/op             216 B/op          2 allocs/op
Benchmark_StdLoggerPrint-4       1061575              1139 ns/op             224 B/op          3 allocs/op
Benchmark_Print-4                1141396              1052 ns/op             216 B/op          2 allocs/op
Benchmark_StdPrint-4             1000000              1134 ns/op             224 B/op          3 allocs/op
Benchmark_Printf-4                913729              1163 ns/op             224 B/op          3 allocs/op
Benchmark_StdPrintf-4             940375              1220 ns/op             240 B/op          3 allocs/op
Benchmark_Sprint-4               1000000              1152 ns/op             224 B/op          3 allocs/op
Benchmark_StdSprint-4            1000000              1184 ns/op             232 B/op          3 allocs/op
PASS
ok      github.com/qq51529210/log       12.809s
```

