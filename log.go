package log

import "sync"

var (
	logPool = sync.Pool{}
)

func init() {
	logPool.New = func() interface{} {
		return new(Log)
	}
}

// A buffer used for format a line of logs.
type Log struct {
	// Log buffer.
	line []byte
	// Integer format buffer.
	fmt []byte
}

// Implements io.Writer interface.
func (l *Log) Write(b []byte) (int, error) {
	l.line = append(l.line, b...)
	return len(b), nil
}

// Return buffer.
func (l *Log) Data() []byte {
	return l.line
}

// Return string of buffer.
func (l *Log) String() string {
	return string(l.line)
}

// Reset buffer.
func (l *Log) Reset() {
	l.line = l.line[:0]
}

// Write a integer into buffer with right align format.
// If len(integer) < length, add 0 to the left.
// Example: 12 -> 0012 while length=4.
func (l *Log) WriteRightAlignInt(integer, length int) {
	// Zero.
	if integer == 0 {
		for i := 0; i < length; i++ {
			l.line = append(l.line, '0')
		}
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.fmt = l.fmt[:0]
	for integer > 0 {
		l.fmt = append(l.fmt, byte('0'+integer%10))
		integer /= 10
	}
	// Add 0 to the left if len(integer) < length
	if length > len(l.fmt) {
		for i := len(l.fmt); i < length; i++ {
			l.line = append(l.line, '0')
		}
	}
	// buff[4,3,2,1]->line[1,2,3,4]
	for i := len(l.fmt) - 1; i >= 0; i-- {
		l.line = append(l.line, l.fmt[i])
	}
}

// Write a integer into buffer with left align format.
// If len(integer) < length,add 0 to the right.
// Example:
//	12 -> 1200 while length=4.
//	1234 -> 12 while length=2.
func (l *Log) WriteLeftAlignInt(integer, length int) {
	// Zero.
	if integer == 0 {
		for i := 0; i < length; i++ {
			l.line = append(l.line, '0')
		}
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.fmt = l.fmt[:0]
	for integer > 0 {
		l.fmt = append(l.fmt, byte('0'+integer%10))
		integer /= 10
	}
	if length < len(l.fmt) {
		// buff[4,3,2,1]->line[1,2]
		for i := len(l.fmt) - 1; i >= len(l.fmt)-length; i-- {
			l.line = append(l.line, l.fmt[i])
		}
	} else {
		// buff[4,3,2,1]->line[1,2,3,4]
		for i := len(l.fmt) - 1; i >= 0; i-- {
			l.line = append(l.line, l.fmt[i])
		}
		// Add 0 to the right.
		for i := len(l.fmt); i < length; i++ {
			l.line = append(l.line, '0')
		}
	}
}

// Write a integer into buffer without algin format.
func (l *Log) WriteInt(integer int) {
	// Zero.
	if integer == 0 {
		l.line = append(l.line, '0')
		return
	}
	// Negative to positive.
	if integer < 0 {
		// Minus sign
		l.line = append(l.line, '-')
		integer = -integer
	}
	// 1234 -> buff[4,3,2,1]
	l.fmt = l.fmt[:0]
	for integer > 0 {
		l.fmt = append(l.fmt, byte('0'+integer%10))
		integer /= 10
	}
	// buff[4,3,2,1]->line[1,2,3,4]
	for i := len(l.fmt) - 1; i >= 0; i-- {
		l.line = append(l.line, l.fmt[i])
	}
}

// Write a byte into buffer.
func (l *Log) WriteUint8(c byte) {
	l.line = append(l.line, c)
}

// Write binary array into buffer.
func (l *Log) WriteBytes(b []byte) {
	l.line = append(l.line, b...)
}

// Write a string into buffer.
func (l *Log) WriteString(s string) {
	l.line = append(l.line, s...)
}
