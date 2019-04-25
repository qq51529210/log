# log
<h3>日志库，又造了个轮子</h3>
<p>
标准库的log包输出格式固定了，所以我把它拆分开了，可以自由的打印日志头部哦。
<br>
Log结构表示一行日志，有自己的缓存，输出到write中是整条日志的缓存，日志不会乱。
<br>
使用了sync.Pool来缓存Log对象，并发安全的。
</p>
<p>
还实现了几个io.Writer<br>
<ul>
<li>LoggerFile，保存到本地磁盘文件</li>
<li>LoggerKafka，保存到Kafka，未实现</li>
<li>LoggerFlume，保存到Flume，未实现</li>
</ul>
</p>

<h3>使用方法，在log_test.go文件中</h3>
<pre>
func TestLog(t *testing.T) {
	level := []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelPanic,
	}
	// 打印，级别，堆栈，调用方法就行，并发安全的。
	// 这几个函数是有换行的
	for i := 0; i < 5; i++ {
		Print(os.Stderr, level[i], i, FileLine(i%3), str)
		Printf(os.Stderr, level[i], i, FileLine(i%3), "Printf: %s", str)
		Sprint(os.Stderr, LevelDebug, i, FileLineName, "Sprint: ", str)
	}
	// 不满意我实现的Print，可以自己定义格式，但是换行要自己打印
	logger := Get()
	// 从缓存池里拿出来的，先清空原来的缓存
	logger.Reset()
	// 1.打印级别
	logger.Level(LevelPanic)
	// 2.打印调用堆栈
	logger.FilePathLine(0, FileLineName)
	// 3.打印时间，格式是
	logger.DateTime(6)
	// 4.打印文本
	logger.String(str)
	// 换行
	logger.EndLine()
	// 输出
	os.Stderr.Write(logger.Bytes())
	// 封装的Print方法
	logger.D(os.Stderr, FileLineFullPath, str)
	logger.I(os.Stderr, FileLineFullPath, str)
	logger.W(os.Stderr, FileLineFullPath, str)
	logger.E(os.Stderr, FileLineFullPath, str)
	Put(logger)
	// 设置自己的分隔符
	// 日期格式 2006-01-02
	// 时间格式 15:04:05.999999999
	DateSeparator = '#'
	TimeSeparator = '*'
	NanoSecSeparator = '>'
	SpaceSeparator = '_'
	FileLineSeparator = '|'
	Print(os.Stderr, LevelDebug, 0, FileLineFullPath, str)
}

// 只打印panic的行
func TestPanic(t *testing.T) {
	defer Recover(os.Stderr,nil)
	Panic("test panic")
}

// 打印完整的堆栈的行
func TestPanicStd(t *testing.T) {
	defer Recover(os.Stderr,nil)
	panic("test std panic")
}

</pre>

<h3>测试</h3>
<p>
既然是造轮子，肯定得造好一点的。
<br>
下一步实现一个简单的sprintf，实现0分配内存，并提升性能。
</p>
<pre>
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
Benchmark_StdLog-4               1000000              1209 ns/op              32 B/op          2 allocs/op
Benchmark_MyLog-4                1000000              1138 ns/op               0 B/op          0 allocs/op
Benchmark_Fmt_StdLog-4           1000000              1249 ns/op              48 B/op          2 allocs/op
Benchmark_Fmt_MyLog-4            1000000              1275 ns/op              16 B/op          1 allocs/op
Benchmark_Log-4                  1000000              1033 ns/op               0 B/op          0 allocs/op
Benchmark_fmt_Log-4              1000000              1189 ns/op              16 B/op          1 allocs/op
Benchmark_Panic-4                2000000               951 ns/op               0 B/op          0 allocs/op
Benchmark_PanicStd-4               10000            107119 ns/op              16 B/op          1 allocs/op
</pre>
