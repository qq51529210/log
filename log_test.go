package log

import "testing"

func Test_writeInt(t *testing.T) {
	l := &Log{}
	// 整数
	n := 0
	l.WriteInt(n)
	if string(l.b) != "0" {
		t.FailNow()
	}
	// 整数
	n = 12345
	l.Reset()
	l.WriteInt(n)
	if string(l.b) != "12345" {
		t.FailNow()
	}
	// 负数
	n = -n
	l.Reset()
	l.WriteInt(n)
	if string(l.b) != "-12345" {
		t.FailNow()
	}
}

func Test_writeIntLeftAlign(t *testing.T) {
	l := &Log{}
	// 整数
	n := 0
	l.WriteIntLeftAlign(n, 5)
	if string(l.b) != "00000" {
		t.FailNow()
	}
	// 整数
	n = 123
	l.Reset()
	l.WriteIntLeftAlign(n, 6)
	if string(l.b) != "123000" {
		t.FailNow()
	}
	// 负数
	n = -n
	l.Reset()
	l.WriteIntLeftAlign(n, 5)
	if string(l.b) != "-12300" {
		t.FailNow()
	}
}

func Test_writeIntRightAlign(t *testing.T) {
	l := &Log{}
	// 整数
	n := 0
	l.WriteIntRightAlign(n, 3)
	if string(l.b) != "000" {
		t.FailNow()
	}
	// 整数
	n = 123
	l.Reset()
	l.WriteIntRightAlign(n, 4)
	if string(l.b) != "0123" {
		t.FailNow()
	}
	// 负数
	n = -n
	l.Reset()
	l.WriteIntRightAlign(n, 5)
	if string(l.b) != "-00123" {
		t.FailNow()
	}
}
