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

type StackInfo byte

const (
	StackInfoDisable StackInfo = iota // 不打印堆栈信息
	StackInfoFile                     // 打印文件名称
	StackInfoPath                     // 打印文件的完整路径
)

// year-month-day hour:minute:second.nano，日期时间格式，分隔符可替换
var (
	logPool                     = sync.Pool{}                 // Logger缓存
	SpaceSeparator    byte      = ' '                         // 空格
	DateSeparator     byte      = '-'                         // 日期
	TimeSeparator     byte      = ':'                         // 时间
	NanoSecSeparator  byte      = '.'                         // 纳秒
	FileLineSeparator byte      = ':'                         // 堆栈
	NanoSecLength               = 6                           // 打印纳秒的长度
	panicFileLine               = []byte("/runtime/panic.go") // 获取panic堆栈判断（只获取panic那行）
	defaultWriter     io.Writer = os.Stdout                   // 默认输出
	defaultStack                = StackInfoPath               // 默认堆栈
)

// 缓存池的new函数
func init() {
	logPool.New = func() interface{} {
		return &Logger{}
	}
}

// 表示一行日志
type Logger struct {
	b []byte // 缓存
}

// io.WriteTo接口
func (lg *Logger) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(lg.b)
	if err != nil {
		return 0, err
	}
	lg.Reset()
	return int64(n), err
}

// io.Writer接口
func (lg *Logger) Write(b []byte) (int, error) {
	lg.b = append(lg.b, b...)
	return len(b), nil
}

// 重置缓存
func (lg *Logger) Reset() {
	lg.b = lg.b[:0]
}

// 返回缓存数据
func (lg *Logger) Data() []byte {
	return lg.b
}

// v字符串长度不满n，右边补0。例，v：12，n:4，1200
func (lg *Logger) WriteIntR0(v, n int) {
	// 值是倒转的
	i1 := len(lg.b)
	if v < 0 {
		lg.b = append(lg.b, '-')
		v = -v
		i1++
	}
	end := i1 + n
	for {
		lg.b = append(lg.b, byte('0'+v%10))
		v /= 10
		if v == 0 {
			break
		}
		n--
	}
	// 先反转
	c := byte(0)
	i2 := len(lg.b) - 1
	for i1 < i2 {
		c = lg.b[i1]
		lg.b[i1] = lg.b[i2]
		lg.b[i2] = c
		i2--
		i1++
	}
	// 再后面补0
	for n > 0 {
		lg.b = append(lg.b, byte('0'))
		n--
	}
	if len(lg.b) > end {
		lg.b = lg.b[:end]
	}
}

// v字符串长度不满n，左边补0。例，v：12，n:4，0012
func (lg *Logger) WriteIntL0(v, n int) {
	// 值是倒转的
	i1 := len(lg.b)
	if v < 0 {
		lg.b = append(lg.b, '-')
		v = -v
		i1++
	}
	for {
		lg.b = append(lg.b, byte('0'+v%10))
		v /= 10
		if v == 0 {
			break
		}
		n--
	}
	// 继续在后面补0
	for n > 1 {
		lg.b = append(lg.b, byte('0'))
		n--
	}
	// 反转
	i2 := len(lg.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = lg.b[i1]
		lg.b[i1] = lg.b[i2]
		lg.b[i2] = c
		i2--
		i1++
	}
}

// 写入一个整数
func (lg *Logger) WriteInt(n int) {
	i1 := len(lg.b)
	if n < 0 {
		lg.b = append(lg.b, '-')
		n = -n
		i1++
	}
	for {
		lg.b = append(lg.b, byte('0'+n%10))
		n /= 10
		if n == 0 {
			break
		}
	}
	i2 := len(lg.b) - 1
	c := byte(0)
	for i1 < i2 {
		c = lg.b[i1]
		lg.b[i1] = lg.b[i2]
		lg.b[i2] = c
		i2--
		i1++
	}
}

// 写入一个字符c
func (lg *Logger) WriteByte(c byte) {
	lg.b = append(lg.b, c)
}

// 写入换行'\n'
func (lg *Logger) WriteEndLine() {
	lg.b = append(lg.b, '\n')
}

// 写入空格' '
func (lg *Logger) WriteSpace() {
	lg.b = append(lg.b, SpaceSeparator)
}

// 写入字符串
func (lg *Logger) WriteString(s string) {
	lg.b = append(lg.b, s...)
}

// 写入数据
func (lg *Logger) WriteBytes(b []byte) {
	lg.b = append(lg.b, b...)
}

// 写入日期时间（year-month-day hour:minute:second.nano）
func (lg *Logger) WriteDateTime(t *time.Time) {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	lg.WriteIntL0(year, 4)
	lg.b = append(lg.b, DateSeparator)
	lg.WriteIntL0(int(month), 2)
	lg.b = append(lg.b, DateSeparator)
	lg.WriteIntL0(day, 2)
	lg.b = append(lg.b, SpaceSeparator)
	lg.WriteIntL0(hour, 2)
	lg.b = append(lg.b, TimeSeparator)
	lg.WriteIntL0(minute, 2)
	lg.b = append(lg.b, TimeSeparator)
	lg.WriteIntL0(second, 2)
	if NanoSecLength > 0 {
		lg.b = append(lg.b, NanoSecSeparator)
		lg.WriteIntR0(t.Nanosecond(), NanoSecLength)
	}
}

// 写入日志级别
func (lg *Logger) WriteLevel(l Level) {
	lg.b = append(lg.b, byte(l))
}

// 写入堆栈的文件名
func (lg *Logger) WriteStackFile(f string, ln int) {
	i := len(f) - 1
	for ; i >= 0; i-- {
		if os.IsPathSeparator(f[i]) {
			i++
			break
		}
	}
	lg.b = append(lg.b, f[i:]...)
	lg.b = append(lg.b, FileLineSeparator)
	lg.WriteInt(ln)
}

// 写入堆栈文件的完整路径
func (lg *Logger) WriteStackPath(f string, ln int) {
	lg.b = append(lg.b, f...)
	lg.b = append(lg.b, FileLineSeparator)
	lg.WriteInt(ln)
}

// w:输出，l:日志级别，i:堆栈，s:日志文本，c:堆栈调用层级
func (lg *Logger) Print(w io.Writer, l Level, i StackInfo, c int, s string) (int, error) {
	lg.printHeader(l, i, c)
	lg.WriteString(s)
	lg.WriteEndLine()
	return w.Write(lg.b)
}

// w:输出，l:日志级别，i:堆栈，s:日志文本，c:堆栈调用层级
func (lg *Logger) PrintBytes(w io.Writer, l Level, i StackInfo, c int, b []byte) (int, error) {
	lg.printHeader(l, i, c)
	lg.WriteBytes(b)
	lg.WriteEndLine()
	return w.Write(lg.b)
}

// w:输出，l:日志级别，i:堆栈，c:堆栈调用层级，f:格式化字符串，a:格式化数据
func (lg *Logger) Printf(w io.Writer, l Level, i StackInfo, c int, f string, a ...interface{}) (int, error) {
	lg.printHeader(l, i, c)
	_, _ = fmt.Fprintf(lg, f, a...)
	lg.WriteEndLine()
	return w.Write(lg.b)
}

// l:日志级别，i:堆栈，c:堆栈调用层级
func (lg *Logger) printHeader(l Level, i StackInfo, c int) {
	lg.Reset()
	lg.WriteLevel(l)
	lg.WriteSpace()
	t := time.Now()
	lg.WriteDateTime(&t)
	lg.WriteSpace()
	switch i {
	case StackInfoDisable:
	case StackInfoFile:
		lg.WriteStackFile(stack(c + 2))
		lg.WriteSpace()
	case StackInfoPath:
		lg.WriteStackPath(stack(c + 2))
		lg.WriteSpace()
	}
}

// w:输出，l:日志级别，i:堆栈，c:堆栈调用层级，a:格式化数据
func (lg *Logger) Sprint(w io.Writer, l Level, i StackInfo, c int, a ...interface{}) (int, error) {
	lg.printHeader(l, i, c)
	_, _ = fmt.Fprint(lg, a...)
	lg.WriteEndLine()
	return w.Write(lg.b)
}

// 写入堆栈信息，panicLine:/runtime/panic.go之后的所有行或者一行
func (lg *Logger) WriteStack(panicLine bool) {
	i1 := len(lg.b)
	i2 := 0
	n := 0
	for {
		lg.b = lg.b[:cap(lg.b)]
		n = runtime.Stack(lg.b[i1:], true)
		i2 = i1 + n
		if i2 < len(lg.b) {
			lg.b = lg.b[:i2]
			break
		}
		lg.b = append(lg.b, make([]byte, 128)...)
	}
	/*
		goroutine 1 [running]:
		main.checkError(...)
		        /Users/ben/Documents/project/go/src/test/main.go:25
		main.main()
		        /Users/ben/Documents/project/go/src/test/main.go:20 +0xb7
	*/
	// 简化一下，只保留文件路径
	n = i2 - 1
	i := i1
	m := i
	// 是否找到/runtime/panic.go，下一行就是panic的地方
	ok := false
Loop:
	for i < n {
		// 文件行开始，'\t'
		if lg.b[i] == '\t' {
			i++
			i1 = i
			for ; i < n; i++ {
				// 文件行路径结束
				if lg.b[i] == ' ' || lg.b[i] == '\n' {
					// 找到/runtime/panic.go
					if !ok {
						ok = bytes.Contains(lg.b[i1:i], panicFileLine)
					} else {
						m += copy(lg.b[m:], lg.b[i1:i])
						if !panicLine {
							// 记录所有行，空格分开
							lg.b[m] = ' '
							m++
						} else {
							// 只记录一行，退出所有，返回
							break Loop
						}
					}
					// 退出这一行，开始下一行
					break
				}
			}
			continue
		}
		i++
	}
	lg.b = lg.b[:m]
	return
}

// 打印输出到默认Writer，l:日志级别，i:堆栈，s:日志文本，
func Print(l Level, c int, s string) (int, error) {
	lg := logPool.Get().(*Logger)
	lg.Reset()
	n, e := lg.Print(defaultWriter, l, defaultStack, c+1, s)
	logPool.Put(lg)
	return n, e
}

// 格式化输出，使用默认的writer和stackinfo
func Printf(l Level, c int, f string, a ...interface{}) (int, error) {
	lg := logPool.Get().(*Logger)
	lg.Reset()
	n, e := lg.Printf(defaultWriter, l, defaultStack, c+1, f, a...)
	logPool.Put(lg)
	return n, e
}

// 格式化输出，使用默认的writer和stackinfo
func Sprint(l Level, c int, a ...interface{}) (int, error) {
	lg := logPool.Get().(*Logger)
	lg.Reset()
	n, e := lg.Sprint(defaultWriter, l, defaultStack, c+1, a...)
	logPool.Put(lg)
	return n, e
}

func stack(n int) (string, int) {
	_, f, l, o := runtime.Caller(n + 1)
	if !o {
		return "???", -1
	}
	return f, l
}

// 简洁的Debug，使用默认的writer和stackinfo
func Debug(a ...interface{}) {
	_, _ = Sprint(LevelDebug, 1, a...)
}

// 简洁的Info，使用默认的writer和stackinfo
func Info(a ...interface{}) {
	_, _ = Sprint(LevelInfo, 1, a...)
}

// 简洁的Warn，使用默认的writer和stackinfo
func Warn(a ...interface{}) {
	_, _ = Sprint(LevelWarn, 1, a...)
}

// 简洁的Error，使用默认的writer和stackinfo
func Error(a ...interface{}) {
	_, _ = Sprint(LevelError, 1, a...)
}

// 简洁的Debug，使用默认的writer和stackinfo
func DebugStack(c int, a ...interface{}) {
	_, _ = Sprint(LevelDebug, c+1, a...)
}

// 简洁的Info，使用默认的writer和stackinfo
func InfoStack(c int, a ...interface{}) {
	_, _ = Sprint(LevelInfo, c+1, a...)
}

// 简洁的Warn，使用默认的writer和stackinfo
func WarnStack(c int, a ...interface{}) {
	_, _ = Sprint(LevelWarn, c+1, a...)
}

// 简洁的Error，使用默认的writer和stackinfo
func ErrorStack(c int, a ...interface{}) {
	_, _ = Sprint(LevelError, c+1, a...)
}

// 设置默认的输出writer
func SetWriter(w io.Writer) {
	if nil != w {
		defaultWriter = w
	}
}

// 设置默认的堆栈
func SetStackInfo(i StackInfo) {
	defaultStack = i
}

// recover，然后输出堆栈，panicLine:/runtime/panic.go之后的所有行或者一行
func Recover(re interface{}, panicLine bool) {
	// 获取Logger
	lg := logPool.Get().(*Logger)
	lg.Reset()
	lg.WriteLevel(LevelPanic)
	lg.WriteSpace()
	t := time.Now()
	lg.WriteDateTime(&t)
	lg.WriteSpace()
	lg.WriteStack(panicLine)
	lg.WriteSpace()
	_, _ = fmt.Fprint(lg, re)
	lg.WriteEndLine()
	_, _ = defaultWriter.Write(lg.b)
	logPool.Put(lg)
}

// 检查错误，直接panic
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
