package log

import "testing"

func Test_writeInt(t *testing.T) {
	b := &Buffer{}
	// 整数
	n := 0
	b.WriteInt(n)
	if string(b.b) != "0" {
		t.FailNow()
	}
	// 整数
	n = 12345
	b.Reset()
	b.WriteInt(n)
	if string(b.b) != "12345" {
		t.FailNow()
	}
	// 负数
	n = -n
	b.Reset()
	b.WriteInt(n)
	if string(b.b) != "-12345" {
		t.FailNow()
	}
}

func Test_writeIntLeftAlign(t *testing.T) {
	b := &Buffer{}
	// 整数
	n := 0
	b.WriteIntLeftAlign(n, 5)
	if string(b.b) != "00000" {
		t.FailNow()
	}
	// 整数
	n = 123
	b.Reset()
	b.WriteIntLeftAlign(n, 6)
	if string(b.b) != "123000" {
		t.FailNow()
	}
	// 负数
	n = -n
	b.Reset()
	b.WriteIntLeftAlign(n, 5)
	if string(b.b) != "-12300" {
		t.FailNow()
	}
}

func Test_writeIntRightAlign(t *testing.T) {
	b := &Buffer{}
	// 整数
	n := 0
	b.WriteIntRightAlign(n, 3)
	if string(b.b) != "000" {
		t.FailNow()
	}
	// 整数
	n = 123
	b.Reset()
	b.WriteIntRightAlign(n, 4)
	if string(b.b) != "0123" {
		t.FailNow()
	}
	// 负数
	n = -n
	b.Reset()
	b.WriteIntRightAlign(n, 5)
	if string(b.b) != "-00123" {
		t.FailNow()
	}
}
