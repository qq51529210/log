package log

import (
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	// Check out std output.
	logger := NewLogger(os.Stderr)
	logger.Print("test 1")
	logger.Printf("test %d", 2)
	logger.Fprint("test", 3)
}
