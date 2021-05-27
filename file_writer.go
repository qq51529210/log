package log

import (
	"os"
	"sync"
	"time"
)

// import (
// 	"bytes"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"sync"
// 	"time"
// )

// const (
// 	day = time.Hour * 24
// )

// var (
// 	errFileLoggerClosed = errors.New("file logger has been closed")
// )

// // 配置，文件日志
// type FileConfig struct {
// 	Dir        string        `json:"dir"`         // 根目录，默认log
// 	Size       int           `json:"size"`        // 每个文件的大小，默认1m
// 	Day        int           `json:"day"`         // 保存的天数，默认1
// 	DayFormat  string        `json:"day_format"`  // 日期目录命名规则
// 	FileFormat string        `json:"file_format"` // 日期目录下文件的命名规则
// 	Duration   time.Duration `json:"duration"`    // 保存到磁盘的间隔，默认1s
// }

// It is a io.Writer to receive data and save data to local disk file.
// First it saves the data in memory and outputs it to a file in next output time.
// If one data file size bigger than File.maxSize,it will create a new one.
// Auto delete files that have been saved for more than File.day days.
// The directory structure is
// root/date.../time...
type File struct {
	mutex sync.Mutex
	// Root directory
	dir     string
	maxSize int64
	maxDay  int
	// Current data.
	data []byte
	// The interval of flush data to file.
	dur time.Duration
	// Directory name format,use Time.Format(),only format date
	dirFormat string
	// File name format,use Time.Format(),only format time
	fileFormat string
	// The file is opened to write.
	file *os.File
}

// // 新的日志文件对象
// func NewFileLogger(cfg *FileConfig) (*File, error) {
// 	if cfg.Day < 1 {
// 		cfg.Day = 1
// 	}
// 	if cfg.Size < 1 {
// 		cfg.Size = 1024 * 1024
// 	}
// 	if cfg.Duration < 1 {
// 		cfg.Duration = 1000
// 	}
// 	if cfg.DayFormat == "" {
// 		cfg.DayFormat = "20060102"
// 	}
// 	if cfg.FileFormat == "" {
// 		cfg.FileFormat = "150405.999999999"
// 	}
// 	if cfg.Dir == "" {
// 		cfg.Dir = "log"
// 	}
// 	// 对象
// 	f := &File{
// 		rootDir:    cfg.Dir,
// 		exit:       make(chan struct{}),
// 		maxDay:     cfg.Day,
// 		maxSize:    cfg.Size,
// 		duration:   cfg.Duration * time.Millisecond,
// 		dayFormat:  cfg.DayFormat,
// 		fileFormat: cfg.FileFormat,
// 	}
// 	f.newFile()
// 	// 保存routine
// 	go func(f *File) {
// 		defer Recover(func() {
// 			f.syncTimer.Stop()
// 			close(f.exit)
// 		})
// 		f.syncTimer = time.NewTimer(f.duration)
// 		for !f.closed {
// 			<-f.syncTimer.C
// 			// 保存数据到文件
// 			f.mutex.Lock()
// 			if f.buffer.Len() > 0 {
// 				// 数据写入文件
// 				n, err := io.Copy(f.file, &f.buffer)
// 				f.printError(err)
// 				f.curSize += n
// 				// 无论写入成功与否，清空缓存
// 				f.buffer.Reset()
// 				// 写入的数据到达最大，开始新文件
// 				if f.curSize >= int64(f.maxSize) {
// 					f.newFile()
// 				}
// 			}
// 			f.mutex.Unlock()
// 			// 检查过期的日志
// 			f.checkExpired()
// 			// 重新计时
// 			f.syncTimer.Reset(f.duration)
// 		}
// 	}(f)
// 	return f, nil
// }

// Implements io.Writer interface.
// func (f *File) Write(b []byte) (int, error) {
// f.mutex.Lock()
// n, e := f.buffer.Write(b)
// f.mutex.Unlock()
// return n, e
// }

// // 清理过期的文件
// func (f *File) checkExpired() {
// 	// 检查过期的
// 	fs, err := ioutil.ReadDir(f.rootDir)
// 	if nil != err {
// 		// 目录不存在，创建
// 		if os.IsNotExist(err) {
// 			err = os.MkdirAll(f.rootDir, os.ModePerm)
// 			if nil == err {
// 				return
// 			}
// 		}
// 		_, _ = os.Stderr.WriteString(err.Error())
// 		return
// 	}
// 	t := time.Now().Add(-day)
// 	for i := 0; i < len(fs); i++ {
// 		if fs[i].ModTime().Sub(t) < 0 {
// 			f.printError(os.RemoveAll(filepath.Join(f.rootDir, fs[i].Name())))
// 		}
// 	}
// }

// // 将缓存的数据保存到磁盘文件中
// func (f *File) newFile() {
// 	// 关闭旧文件
// 	f.closeFile()
// 	// 时间
// 	now := time.Now()
// 	// 创建日期目录
// 	dateDir := filepath.Join(f.rootDir, now.Format(f.dayFormat))
// 	err := os.MkdirAll(dateDir, os.ModePerm)
// 	if nil != err {
// 		_, _ = fmt.Fprintln(os.Stderr, err)
// 		return
// 	}
// 	// 新的日志文件
// 	timeFile := filepath.Join(dateDir, now.Format(f.fileFormat))
// 	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
// 	if nil != err {
// 		_, _ = fmt.Fprintln(os.Stderr, err)
// 		return
// 	}
// }

// // 关闭文件
// func (f *File) closeFile() {
// 	if nil != f.file {
// 		// 把剩下的缓存写入文件，再关闭
// 		_, err := io.Copy(f.file, &f.buffer)
// 		f.printError(err)
// 		f.printError(f.file.Close())
// 		f.file = nil
// 	}
// }

// func (f *File) Close() error {
// 	f.mutex.Lock()
// 	if f.closed {
// 		f.mutex.Unlock()
// 		return errFileLoggerClosed
// 	}
// 	f.closed = true
// 	f.mutex.Unlock()
// 	// 退出
// 	f.syncTimer.Reset(0)
// 	// 等待routine退出
// 	<-f.exit
// 	// 关闭文件
// 	f.closeFile()
// 	return nil
// }

// func (f *File) printError(err error) {
// 	if err != nil {
// 		_, _ = fmt.Fprintln(os.Stderr, err)
// 	}
// }
