# log
标准库的log包输出日志头部格式固定了，我不喜欢，所以我把它拆分开了，可以自由的打印日志头部  
## 用法
如果使用默认的Print函数，输出的格式为<year-month-day hour:minute:second.nano level stack:log>  
进行格式化的是Log{}这个结构体，调用它的函数可以自定义格式化的日志，具体看Print()函数的实现  
## 测试  
既然是造轮子，肯定得造好一点的  
下一步实现一个简单的sprintf，实现0分配内存，并提升性能  
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