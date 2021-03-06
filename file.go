package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	defaultMaxFileSize = 1024 * 1024
)

var (
	errFileClosed = errors.New("file has been closed")
)

// Create a new File output. FlushLoop has been called inside as a goutine.
// Parameters description reference File.FlushLoop().
func NewFileLogger(rootDir string, maxFileSize, maxDay int, dur time.Duration) (*File, error) {
	f := new(File)
	go f.FlushAndCheckLoop(rootDir, maxFileSize, maxDay, dur)
	return f, nil
}

// File implements io.Writer to receive data in memory and  outputs it to a file in next output time.
type File struct {
	lock sync.Mutex
	// Wait for FlushLoop routine exit.
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
	// Log data, waitting for flush.
	data []byte
	// The file is opened to write.
	file *os.File
}

// Parameters description reference File.FlushLoop().
func (f *File) init(rootDir string, maxFileSize, maxKeepDay int, dur time.Duration) {
	f.ok = true
	f.rootDir = rootDir
	f.maxFileSize = maxFileSize
	if f.maxFileSize < 1 {
		f.maxFileSize = defaultMaxFileSize
	}
	f.maxKeepDay = time.Duration(maxKeepDay)
	if f.maxKeepDay < 1 {
		f.maxKeepDay = 1
	}
	f.maxKeepDay *= time.Hour * 24
	f.timer = time.NewTicker(dur)
	f.exit = make(chan struct{})
}

// FlushLoop flush memory data into file and check expired logs.
// rootDir is the root directory of all log files. rootDir has many dateDir, dateDir has many timeFile.
// maxFileSize is the max size of a log file, if output file's size is greater that it, a new log file will be created.
// maxDay is the max days to keep log file.
// dur is the interval to flush data.
func (f *File) FlushAndCheckLoop(rootDir string, maxFileSize, maxDay int, dur time.Duration) {
	f.wait.Add(1)
	defer f.wait.Done()
	// Init data.
	f.lock.Lock()
	if !f.ok {
		f.init(rootDir, maxFileSize, maxDay, dur)
	} else {
		// Routine has been called.
		f.lock.Unlock()
		return
	}
	// Open file.
	f.openFile()
	f.lock.Unlock()
	// Check expire file.
	checkTime := time.Now()
	f.checkExpiredFile(checkTime)
	for f.ok {
		select {
		case now := <-f.timer.C:
			// Time to flush data.
			f.lock.Lock()
			f.flushData()
			f.lock.Unlock()
			// Another day.
			if now.Day() != checkTime.Day() {
				f.checkExpiredFile(now)
				checkTime = now
			}
		case <-f.exit:
			// Close() called.
		}
	}
}

// Append data b to memory.
func (f *File) Write(b []byte) (int, error) {
	f.lock.Lock()
	if f.ok {
		if len(b)+len(f.data) >= f.maxFileSize {
			// Memory data is greater than maxFileSize
			f.flushData()
			f.closeFile()
			f.openFile()
		}
		f.data = append(f.data, b...)
		f.lock.Unlock()
		return len(b), nil
	}
	f.lock.Unlock()
	return 0, errFileClosed
}

// Change File state and wait for FlushLoop exit.
func (f *File) Close() error {
	f.lock.Lock()
	if !f.ok {
		f.lock.Unlock()
		return errFileClosed
	}
	f.ok = false
	f.lock.Unlock()
	close(f.exit)
	// Waitting for FlushLoop exit.
	f.wait.Wait()
	// Stop timer.
	f.timer.Stop()
	// Flush rest of data.
	f.flushData()
	// Close file.
	f.closeFile()
	return nil
}

// Check and remove expired logs.
func (f *File) checkExpiredFile(now time.Time) {
	files, err := ioutil.ReadDir(f.rootDir)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	delTime := now.Add(-f.maxKeepDay)
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

// Open a new file.
func (f *File) openFile() {
	now := time.Now()
	// Create "date" directory first.
	dateDir := filepath.Join(f.rootDir, now.Format("20060102"))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Create "time" log file
	timeFile := filepath.Join(dateDir, strconv.FormatInt(now.Unix(), 10))
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Close file if is opened.
func (f *File) closeFile() {
	if nil != f.file {
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
