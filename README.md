# log
标准库的log包输出日志头部格式固定了，我不喜欢，所以我把它拆分开了，可以自由的打印日志头部  
## 用法
如果使用默认的Print函数，输出的格式为<level year-month-day hour:minute:second.nano stack log>  
可以包装Logger来实现自己的日志格式
```
lg := new(Logger)
// 级别
lg.WriteLevel(l)
lg.WriteSpace()
// 时间加日期
t := time.Now()
lg.WriteDateTime(&t)
lg.WriteSpace()
// 堆栈文件/完整路径
lg.WriteStackFile()
lg.WriteStackPath()
// 日志文本/数据
lg.WriteString()
lg.WriteBytes()

对时间日期格式不满足，也可以自己包装
lg.WriteInt()
lg.WriteIntR0()
lg.WriteIntL0()
```
## 测试  
下面是和标准库的性能测试
```
goos: darwin
goarch: amd64
pkg: gomod/log
Benchmark_LogPrint-4       	 5118735	       230 ns/op	       0 B/op	       0 allocs/op
Benchmark_Std_LogPrint-4   	 3416056	       342 ns/op	       5 B/op	       1 allocs/op
Benchmark_Std_Print-4      	 3407294	       339 ns/op	       5 B/op	       1 allocs/op
Benchmark_Print-4          	 4731691	       249 ns/op	       0 B/op	       0 allocs/op
Benchmark_Printf-4         	 3255664	       365 ns/op	       8 B/op	       1 allocs/op
Benchmark_Std_Printf-4     	 2807124	       426 ns/op	      24 B/op	       2 allocs/op
Benchmark_Sprint-4         	 3290929	       357 ns/op	       8 B/op	       1 allocs/op
Benchmark_Std_Sprint-4     	 3159760	       378 ns/op	      16 B/op	       2 allocs/op
PASS
```
