package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Create a new File output.
func NewFileLogger(dir string, maxFileSize, maxDay int, dur time.Duration) (*File, error) {
	f := new(File)
	go f.FlushLoop(dir, maxFileSize, maxDay, dur)
	return f, nil
}

// It is a io.Writer to receive data and save data to local disk file.
// First it saves the data in memory and outputs it to a file in next output time.
// If one data file size bigger than File.maxSize,it will create a new one.
// Auto delete files saved before File.day days.
// The directory structure is "root/date/time".
type File struct {
	lock sync.Mutex
	wait sync.WaitGroup
	// Timer for flush loop.
	timer *time.Ticker
	// Close signal.
	exit chan struct{}
	// If closed.
	ok bool
	// Root directory.
	rootDir string
	// Max size of a log file.
	maxFileSize int
	// Max days to keep file.
	maxKeepDay time.Duration
	// Log data.
	data []byte
	// The file is opened to write.
	file *os.File
}

func (f *File) init(rootDir string, maxFileSize, maxKeepDay int, dur time.Duration) {
	f.ok = true
	f.rootDir = rootDir
	f.maxFileSize = maxFileSize
	if f.maxFileSize < 1 {
		f.maxFileSize = 1024 * 1024
	}
	f.maxKeepDay = time.Duration(maxKeepDay)
	if f.maxKeepDay < 1 {
		f.maxKeepDay = 1
	}
	f.maxKeepDay *= -time.Hour * 24
	f.timer = time.NewTicker(dur)
	f.exit = make(chan struct{})
}

// Flush memory data into file loop.
// Arg dur is the flush interval.
// Arg dir is the root of log dir.
func (f *File) FlushLoop(dir string, maxFileSize, maxDay int, dur time.Duration) {
	f.wait.Add(1)
	defer f.wait.Done()
	startTime := time.Now()
	// Init if not
	f.lock.Lock()
	if !f.ok {
		f.init(dir, maxFileSize, maxDay, dur)
	} else {
		// Routine has been called.
		f.lock.Unlock()
		return
	}
	// Check expire file first.
	f.checkExpire(startTime)
	// Open file.
	f.openFile()
	f.lock.Unlock()
	// Loop
	for f.ok {
		select {
		case now := <-f.timer.C:
			// Time to flush data.
			f.lock.Lock()
			f.flushData()
			// Another day.
			if now.Day() != startTime.Day() {
				f.checkExpire(now)
				startTime = now
			}
			f.lock.Unlock()
		case <-f.exit:
			// Close() called.
		}
	}
}

// Implements io.Writer interface.
// Append data b to memory.
func (f *File) Write(b []byte) (int, error) {
	f.lock.Lock()
	f.data = append(f.data, b...)
	if len(f.data) > f.maxFileSize {
		f.flushData()
		f.closeFile()
		f.openFile()
	}
	f.lock.Unlock()
	return len(b), nil
}

// Implements io.Closer interface.
// Flush memory data and close file.
// Wait for flush routine exit.
func (f *File) Close() error {
	f.lock.Lock()
	if !f.ok {
		f.lock.Unlock()
		return errors.New("file has been closed")
	}
	f.ok = false
	f.lock.Unlock()
	// Wait for saving routine exit.
	f.wait.Wait()
	// Stop timer.
	f.timer.Stop()
	// Close file.
	f.closeFile()
	return nil
}

// Check and remove expired directory.
func (f *File) checkExpire(now time.Time) {
	files, err := ioutil.ReadDir(f.rootDir)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	delTime := time.Now().Add(f.maxKeepDay)
	for i := 0; i < len(files); i++ {
		// Remove all file which modify time is before delTime.
		if files[i].ModTime().Sub(delTime) < 0 {
			err = os.RemoveAll(filepath.Join(f.rootDir, files[i].Name()))
			if nil != err {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// Close old file and open a new file.
func (f *File) openFile() {
	now := time.Now()
	// Create directory first.
	dateDir := filepath.Join(f.rootDir, now.Format("20060102"))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Create log file
	timeFile := filepath.Join(dateDir, now.Format("150405"))
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Close file if is opened.
func (f *File) closeFile() {
	if nil != f.file {
		// Rest data in the memory.
		if len(f.data) > 0 {
			f.file.Write(f.data)
			f.data = f.data[:0]
		}
		f.file.Close()
		f.file = nil
	}
}

// Flush data to file.
func (f *File) flushData() {
	_, err := f.file.Write(f.data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	f.data = f.data[:0]
}
