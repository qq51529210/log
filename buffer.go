package log

import "sync"

var (
	// 用于格式化整数
	intByte []byte
	// 缓存池
	bufPool sync.Pool
)

func init() {
	for i := 0; i < 10; i++ {
		intByte = append(intByte, '0'+byte(i))
	}
	//
	bufPool.New = func() any {
		return new(Buffer)
	}
}

// Buffer 用于写入一行日志
type Buffer struct {
	// 缓存
	b []byte
	// 用于整数格式化
	f []byte
}

// Reset 重置缓存
func (b *Buffer) Reset() {
	b.b = b.b[:0]
}

// Write 实现 io.Writer
func (b *Buffer) Write(data []byte) (int, error) {
	b.b = append(b.b, data...)
	return len(data), nil
}

// writeReversedInt 写入反转的整数
func (b *Buffer) writeReversedInt(v int) {
	b.f = b.f[:0]
	// 零
	if v == 0 {
		b.f = append(b.f, '0')
		return
	}
	// 1234 -> buff[4,3,2,1]
	for v > 0 {
		b.f = append(b.f, intByte[v%10])
		v /= 10
	}
}

// reversed 将 b.f 的值反转，然后写入 b
func (b *Buffer) reversed() {
	for i, j := 0, len(b.f)-1; i < j; i, j = i+1, j-1 {
		b.f[i], b.f[j] = b.f[j], b.f[i]
	}
	b.b = append(b.b, b.f...)
}

// WriteInt 写入整数
func (b *Buffer) WriteInt(v int) {
	if v < 0 {
		// 负数
		b.writeReversedInt(-v)
		b.f = append(b.f, '-')
	} else {
		// 正数
		b.writeReversedInt(v)
	}
	b.reversed()
}

// WriteIntRightAlign 写入整数，左侧补齐 n 个 0
func (b *Buffer) WriteIntRightAlign(v, n int) {
	if v < 0 {
		// 负数
		b.writeReversedInt(-v)
		for i := len(b.f); i < n; i++ {
			b.f = append(b.f, '0')
		}
		b.f = append(b.f, '-')
	} else {
		// 正数
		b.writeReversedInt(v)
		for i := len(b.f); i < n; i++ {
			b.f = append(b.f, '0')
		}
	}
	b.reversed()
}

// WriteIntLeftAlign 写入整数，右侧补齐 0
func (b *Buffer) WriteIntLeftAlign(v, n int) {
	if v < 0 {
		// 负数
		b.writeReversedInt(-v)
		n -= len(b.f)
		b.f = append(b.f, '-')
	} else {
		// 正数
		b.writeReversedInt(v)
		n -= len(b.f)
	}
	b.reversed()
	for i := 0; i < n; i++ {
		b.b = append(b.b, '0')
	}
}

// WriteString 写入字符串
func (b *Buffer) WriteString(s string) {
	b.b = append(b.b, s...)
}

// WriteBytes 写入字节数组
func (b *Buffer) WriteBytes(s []byte) {
	b.b = append(b.b, s...)
}

// WriteByte 写入字节
func (b *Buffer) WriteByte(s byte) {
	b.b = append(b.b, s)
}
