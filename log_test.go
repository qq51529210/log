package log

import "testing"

func TestLog_WriteRightAlignInt(t *testing.T) {
	log := new(Log)

	log.WriteRightAlignInt(1234, 7)
	if log.String() != "0001234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(-1234, 3)
	if log.String() != "-1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(-1234, 6)
	if log.String() != "-001234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteRightAlignInt(0, 4)
	if log.String() != "0000" {
		t.FailNow()
	}
}

func TestLog_WriteLeftAlignInt(t *testing.T) {
	log := new(Log)

	log.WriteLeftAlignInt(1234, 7)
	if log.String() != "1234000" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(-1234, 2)
	if log.String() != "-12" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(-1234, 6)
	if log.String() != "-123400" {
		t.FailNow()
	}

	log.Reset()
	log.WriteLeftAlignInt(0, 4)
	if log.String() != "0000" {
		t.FailNow()
	}
}

func TestLog_WriteInt(t *testing.T) {
	log := new(Log)

	log.WriteInt(1234)
	if log.String() != "1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteInt(-1234)
	if log.String() != "-1234" {
		t.FailNow()
	}

	log.Reset()
	log.WriteInt(0)
	if log.String() != "0" {
		t.FailNow()
	}
}
