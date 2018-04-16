package log

import (
	"errors"
	"sync"
	"testing"
)

func TestOutput(t *testing.T) {
	l := Open("", 0, 0, true)
	l.Print(LEVEL_DEBUG, "LEVEL_DEBUG")
	l.Print(LEVEL_WARN, "LEVEL_WARN")
	l.Print(LEVEL_INFO, "LEVEL_INFO")
	l.Print(LEVEL_ERROR, "LEVEL_ERROR")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer func() {
			l.RecoverError(recover())
			wg.Done()
		}()
		Panic(errors.New("test log panic"))
	}()
	go func() {
		defer func() {
			l.RecoverError(recover())
			wg.Done()
		}()
		panic(errors.New("test panic"))
	}()
	wg.Wait()
}
