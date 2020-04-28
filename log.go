package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Level byte

const (
	LevelDebug   Level = 'D'
	LevelInfo    Level = 'I'
	LevelWarn    Level = 'W'
	LevelError   Level = 'E'
	LevelPanic   Level = 'P'
	LevelRecover Level = 'R'
)

type FileLine byte

const (
	FileLineDisable  FileLine = iota // 禁用打印堆栈信息
	FileLineFullPath                 // 打印堆栈信息,
	FileLineName                     // 打印文件名称而不是文件的完整路径
)

// year-month-day hour:minute:second.nano，日期时间格式，分隔符可替换
var (
	logPool                = sync.Pool{}                 // Log{}缓存
	errPool                = sync.Pool{}                 // panicError{}缓存
	SpaceSeparator    byte = ' '                         // 空格
	DateSeparator     byte = '-'                         // 日期
	TimeSeparator     byte = ':'                         // 时间
	NanoSecSeparator  byte = '.'                         // 纳秒
	FileLineSeparator byte = ':'                         // 堆栈
	NanoSecLength          = 6                           // 打印纳秒的长度
	unknownFileLine        = []byte("???:-1:")           // 获取堆栈失败时使用
	panicLine              = []byte("/runtime/panic.go") // 获取panic堆栈判断（只获取panic那行）
)

// 缓存池的new函数
func init() {
	logPool.New = func() interface{} {
		return &Log{}
	}
	errPool.New = func() interface{} {
		return &panicError{}
	}
}

type PrintInfo struct {
	Writer   io.Writer
	Level    Level
	Skip     int
	FileLine FileLine
}

// 表示一行日志
type Log struct {
	b    []byte // 缓存
	Info PrintInfo
}

// 重置缓存
func (this *Log) Reset() {
	this.b = this.b[:0]
}

// 返回缓存
func (this *Log) Bytes() []byte {
	return this.b
}

// io.WriteTo接口
func (this *Log) WriteTo(writer io.Writer) (int64, error) {
	// 写入writer
	n, err := writer.Write(this.b)
	if err != nil {
		return 0, err
	}
	// 重置
	this.Reset()
	return int64(n), err
}

// 整数不满length，右边补0。0012-08-02 04:03:05.123000
func (this *Log) IntegerAlignRight(value, length int) {
	// 值是倒转的
	i1 := len(this.b)
	end := i1 + length
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
		length--
	}
	// 先反转
	c := byte(0)
	i2 := len(this.b) - 1
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
	// 再后面补0
	for length > 0 {
		this.b = append(this.b, byte('0'))
		length--
	}
	if len(this.b) > end {
		this.b = this.b[:end]
	}
}

// 整数不满length，左边补0。0012-08-02 04:03:05
func (this *Log) IntegerAlignLeft(value, length int) {
	// 值是倒转的
	i1 := len(this.b)
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
		length--
	}
	// 继续在后面补0
	for length > 1 {
		this.b = append(this.b, byte('0'))
		length--
	}
	// 反转
	i2 := len(this.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
}

// 写入整数value
func (this *Log) Integer(value int) {
	i1 := len(this.b)
	for {
		this.b = append(this.b, byte('0'+value%10))
		value /= 10
		if value == 0 {
			break
		}
	}
	i2 := len(this.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = this.b[i1]
		this.b[i1] = this.b[i2]
		this.b[i2] = c
		i2--
		i1++
	}
}

// 写入一个字符c
func (this *Log) Byte(c byte) {
	this.b = append(this.b, c)
}

// 写入换行'\n'
func (this *Log) EndLine() {
	this.b = append(this.b, '\n')
}

// 写入日期时间（year-month-day hour:minute:second.nano），nsec表示保留纳秒的位数
func (this *Log) DateTime(nsec int) {
	date_time := time.Now()
	year, month, day := date_time.Date()
	hour, minute, second := date_time.Clock()

	this.IntegerAlignLeft(year, 4)
	this.b = append(this.b, DateSeparator)
	this.IntegerAlignLeft(int(month), 2)
	this.b = append(this.b, DateSeparator)
	this.IntegerAlignLeft(day, 2)
	this.b = append(this.b, SpaceSeparator)
	this.IntegerAlignLeft(hour, 2)
	this.b = append(this.b, TimeSeparator)
	this.IntegerAlignLeft(minute, 2)
	this.b = append(this.b, TimeSeparator)
	this.IntegerAlignLeft(second, 2)
	if nsec > 0 {
		// 再长没有意义
		if nsec > 6 {
			nsec = 6
		}
		this.b = append(this.b, NanoSecSeparator)
		this.IntegerAlignRight(date_time.Nanosecond(), nsec)
	}
	this.b = append(this.b, SpaceSeparator)
}

// 写入日志级别
func (this *Log) Level(level Level) {
	this.b = append(this.b, byte(level))
	this.b = append(this.b, SpaceSeparator)
}

// 写入堆栈信息
func (this *Log) FilePathLine(skip int, fileLine FileLine) {
	if fileLine != FileLineDisable {
		_, f, l, o := runtime.Caller(skip + 1)
		if !o {
			this.b = append(this.b, unknownFileLine...)
		} else {
			if FileLineName == fileLine {
				for i := len(f) - 1; i >= 0; i-- {
					if f[i] == '/' {
						this.b = append(this.b, f[i+1:]...)
						this.b = append(this.b, FileLineSeparator)
						this.Integer(l)
						break
					}
				}
			} else {
				this.b = append(this.b, f...)
				this.b = append(this.b, FileLineSeparator)
				this.Integer(l)
			}
		}
		this.b = append(this.b, FileLineSeparator)
	}
	this.b = append(this.b, SpaceSeparator)
}

// 写入字符串
func (this *Log) String(s string) {
	this.b = append(this.b, s...)
}

// 写入数据
func (this *Log) Write(b []byte) (int, error) {
	this.b = append(this.b, b...)
	return len(b), nil
}

func (this *Log) PrintBytes(log []byte) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(this.Info.Level)
	this.FilePathLine(this.Info.Skip+1, this.Info.FileLine)
	this.b = append(this.b, log...)
	this.EndLine()
	return this.Info.Writer.Write(this.b)
}

func (this *Log) Print(log string) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(this.Info.Level)
	this.FilePathLine(this.Info.Skip+1, this.Info.FileLine)
	this.String(log)
	this.EndLine()
	return this.Info.Writer.Write(this.b)
}

func (this *Log) Printf(format string, a ...interface{}) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(this.Info.Level)
	this.FilePathLine(this.Info.Skip+1, this.Info.FileLine)
	fmt.Fprintf(this, format, a...)
	this.EndLine()
	return this.Info.Writer.Write(this.b)
}

func (this *Log) Sprint(a ...interface{}) (int, error) {
	this.Reset()
	this.DateTime(6)
	this.Level(this.Info.Level)
	this.FilePathLine(this.Info.Skip+1, this.Info.FileLine)
	fmt.Fprint(this, a...)
	this.EndLine()
	return this.Info.Writer.Write(this.b)
}

func (this *Log) D(log string) (int, error) {
	this.Info.Level = LevelDebug
	this.Info.Skip = 1
	return this.Print(log)
}

func (this *Log) I(log string) (int, error) {
	this.Info.Level = LevelInfo
	this.Info.Skip = 1
	return this.Print(log)
}

func (this *Log) W(log string) (int, error) {
	this.Info.Level = LevelWarn
	this.Info.Skip = 1
	return this.Print(log)
}

func (this *Log) E(log string) (int, error) {
	this.Info.Level = LevelError
	this.Info.Skip = 1
	return this.Print(log)
}

// 写入堆栈信息，full表示是否整个堆栈，或者只取panic的那一行
func (this *Log) Stack(full bool) {
	i1 := len(this.b)
	i2 := 0
	n := 0
	for {
		this.b = this.b[:cap(this.b)]
		n = runtime.Stack(this.b[i1:], full)
		i2 = i1 + n
		if i2 < len(this.b) {
			this.b = this.b[:i2]
			break
		}
		this.b = append(this.b, make([]byte, 1024)...)
	}
	panic_line := false
	// 简化一下，只保留文件路径
	n = i2 - 1
	m := i1
Loop:
	for i := i1; i < n; {
		// 找出每一行
		if this.b[i] == '\t' {
			this.b[i] = ' '
			// 路径开始
			i++
			i1 = i
			for ; i < n; i++ {
				if this.b[i] == ' ' {
					// 路径结束
					i2 = i + 1
					// 只要一行
					if !full {
						if !panic_line {
							panic_line = bytes.Contains(this.b[i1:i2], panicLine)
						} else {
							m += copy(this.b[m:], this.b[i1:i2])
							break Loop
						}
					} else {
						if !full {
							// 只要panic之后的路径
							if !panic_line {
								panic_line = bytes.Contains(this.b[i1:i2], panicLine)
							} else {
								m += copy(this.b[m:], this.b[i1:i2])
							}
						} else {
							m += copy(this.b[m:], this.b[i1:i2])
						}
					}
					break
				}
			}
		}
		i++
	}
	this.b = this.b[:m]
}

// 从缓存中获取Log{}
func Get() *Log {
	//p := logPool.Get().(*Log)
	//p.Reset()
	//return p
	return logPool.Get().(*Log)
}

// Log{}返回缓存中
func Put(l *Log) {
	logPool.Put(l)
}

// 打印
func Print(writer io.Writer, level Level, skip int, fileLine FileLine, log string) (int, error) {
	l := logPool.Get().(*Log)
	l.Info.Writer = writer
	l.Info.Skip = skip + 1
	l.Info.Level = level
	l.Info.FileLine = fileLine
	n, e := l.Print(log)
	logPool.Put(l)
	return n, e
}

// 格式化输出
func Printf(writer io.Writer, level Level, skip int, fileLine FileLine, format string, a ...interface{}) (int, error) {
	l := logPool.Get().(*Log)
	l.Info.Writer = writer
	l.Info.Skip = skip + 1
	l.Info.Level = level
	l.Info.FileLine = fileLine
	n, e := l.Printf(format, a...)
	logPool.Put(l)
	return n, e
}

// 格式化输出
func Sprint(writer io.Writer, level Level, skip int, fileLine FileLine, a ...interface{}) (int, error) {
	l := logPool.Get().(*Log)
	l.Info.Writer = writer
	l.Info.Skip = skip + 1
	l.Info.Level = level
	l.Info.FileLine = fileLine
	n, e := l.Sprint(a...)
	logPool.Put(l)
	return n, e
}

// 只输出panic的堆栈的路径，其他堆栈信息不输出
// 函数内部会调用recover()
// 参数
// writer: 输出
// full: 输出所有的堆栈，或者只输出发生panic的那一行
// cb: 输出完后回调
// 返回
// 是否发生的panic
func Recover(writer io.Writer, full bool, cb func()) bool {
	// recover
	o := RecoverValue(writer, full, recover())
	// 回调函数
	if cb != nil {
		cb()
	}
	return o
}

// 只输出panic的堆栈的路径，其他堆栈信息不输出
// 参数
// writer: 输出
// full: 输出所有的堆栈，或者只输出发生panic的那一行
// re: 在函数外调用recover()的值
// 返回
// 是否发生的panic
func RecoverValue(writer io.Writer, full bool, re interface{}) bool {
	if re != nil {
		// 获取Log
		l := logPool.Get().(*Log)
		l.Reset()
		if e, o := re.(*panicError); o {
			// 时间
			l.DateTime(6)
			// 级别，recover
			l.Level(LevelRecover)
			// 如果是log.Error，可以提升性能
			l.b = append(l.b, e.File...)
			l.b = append(l.b, FileLineSeparator)
			l.Integer(e.Line)
			l.b = append(l.b, FileLineSeparator)
			l.b = append(l.b, SpaceSeparator)
			l.String(e.Log)
			errPool.Put(e)
		} else {
			// 时间
			l.DateTime(6)
			// 级别，panic
			l.Level(LevelPanic)
			// 信息
			l.String(fmt.Sprint(re))
			l.Byte(SpaceSeparator)
			// 不是panicError，从堆栈找到panic的行
			l.Stack(full)
		}
		l.EndLine()
		writer.Write(l.b)
		logPool.Put(l)
		return true
	}
	return false
}

type panicError struct {
	File string // 文件路径
	Line int    // 文件行
	Log  string // 信息
}

func (this *panicError) Error() string {
	return this.Log
}

// 直接panic
func Panic(log string) {
	panic(newError(1, log))
}

// 检查错误，直接panic
func CheckError(e error) {
	if e != nil {
		panic(newError(1, e.Error()))
	}
}

// 返回panicError{}，带有当前的堆栈信息，提高性能
func newError(skip int, log string) error {
	// 获取Error
	e := errPool.Get().(*panicError)
	_, f, l, o := runtime.Caller(skip + 1)
	if o {
		e.File = f
		e.Line = l
	} else {
		e.File = "???"
		e.Line = -1
	}
	e.Log = log
	return e
}

// 默认的输出writer
var defaultWriter io.Writer = os.Stdout

// 设置默认的输出writer
func SetDefaultWriter(w io.Writer) {
	if nil != w {
		defaultWriter = w
	}
}

// 简洁的Debug
func Debug(a ...interface{}) {
	Sprint(os.Stderr, LevelDebug, 1, FileLineFullPath, a...)
}

// 简洁的Info
func Info(a ...interface{}) {
	Sprint(os.Stderr, LevelInfo, 1, FileLineFullPath, a...)
}

// 简洁的Warn
func Warn(a ...interface{}) {
	Sprint(os.Stderr, LevelWarn, 1, FileLineFullPath, a...)
}

// 简洁的Error
func Error(a ...interface{}) {
	Sprint(os.Stderr, LevelError, 1, FileLineFullPath, a...)
}
