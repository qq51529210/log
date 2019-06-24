package log

import (
	"bytes"
	"container/list"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/qq51529210/common"
)

// 配置
type LoggerFileConfig struct {
	Dir        string `json:"dir"`         // 根目录
	Size       string `json:"size"`        // 每个文件的大小
	Day        int    `json:"day"`         // 保存的天数
	DayFormat  string `json:"day_format"`  // 日期目录命名规则
	FileFormat string `json:"file_format"` // 日期目录下文件的命名规则
	Duration   string `json:"duration"`    // 保存到磁盘的间隔
	Std        string `json:"std"`         // 是否输出到标准流，空或不对，都不输出
}

// 磁盘日志
// 目录结构，/root/date/time.nanosec.log
// 当文件大小到达指定的值，就重新写一个新的文件
// 自动删除超过指定天数的文件夹
// 级别低于指定值不会输出到文件，但是会输出到std流（需配置）
// 只会打开一个文件
type LoggerFile struct {
	mutex      sync.Mutex    // 同步锁
	rootDir    string        // 根目录
	curSize    int64         // 文件写入数据大小
	maxSize    int           // 文件最大大小
	maxDay     int           // 最大天数
	buffer     bytes.Buffer  // 缓存
	level      Level         // 级别
	duration   time.Duration // 保存间隔
	dayFormat  string        // 日期目录的格式
	fileFormat string        // 文件名字格式
	stdout     io.Writer     // 是否输出到stdout
	exit       chan struct{} // 退出
	syncTimer  *time.Timer   // 保存数据的计时器
	file       *os.File      // 文件
	closed     bool          // 已经关闭？
}

// 写数据到缓存
func (this *LoggerFile) Write(b []byte) (int, error) {
	// 输出控制台
	if nil != this.stdout {
		this.stdout.Write(b)
	}
	// 写入缓存
	this.mutex.Lock()
	n, e := this.buffer.Write(b)
	this.mutex.Unlock()
	return n, e
}

// 设置日志的级别
func (this *LoggerFile) SetStdWriter(std string) {
	switch strings.ToLower(std) {
	case "out":
		this.stdout = os.Stdout
	case "err":
		this.stdout = os.Stderr
	default:
		this.stdout = nil
	}
}

// 将缓存的数据保存到磁盘文件中
// 注意，与命名格式无关的全部会被删除
func (this *LoggerFile) checkExpired() {
	// 检查过期的
	fs, e := ioutil.ReadDir(this.rootDir)
	if nil != e {
		if os.IsNotExist(e) {
			e = os.MkdirAll(this.rootDir, os.ModePerm)
			if nil != e {
				Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
			}
			return
		}
		Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
		return
	}
	// 目录数量少
	count := len(fs)
	if count <= this.maxDay {
		return
	}
	// 检查目录，排序
	log_dir := list.New()
	for i := 0; i < len(fs); i++ {
		// 不是目录
		if !fs[i].IsDir() {
			count--
			e = os.RemoveAll(filepath.Join(this.rootDir, fs[i].Name()))
			if nil != e {
				Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
			}
		}
		t, e := time.Parse(this.dayFormat, fs[i].Name())
		// 解析失败，不是日志格式，其他目录
		if nil != e {
			count--
			e = os.RemoveAll(filepath.Join(this.rootDir, fs[i].Name()))
			if nil != e {
				Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
			}
		}
		// 是日志目录
		if log_dir.Len() <= 0 {
			log_dir.PushBack(&t)
		} else {
			for ele := log_dir.Back(); ele != nil; ele = ele.Prev() {
				tt := ele.Value.(*time.Time)
				if t.Sub(*tt) > 0 {
					log_dir.InsertAfter(&t, ele)
				}
			}
		}
	}
	// 删除
	for log_dir.Len() > count {
		ele := log_dir.Front()
		log_dir.Remove(ele)
		t := ele.Value.(*time.Time)
		e = os.RemoveAll(filepath.Join(this.rootDir, t.Format(this.dayFormat)))
		if nil != e {
			Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
		}
	}
}

// 将缓存的数据保存到磁盘文件中
func (this *LoggerFile) saveFile() {
	// 是否有数据
	this.mutex.Lock()
	if this.buffer.Len() > 0 {
		n, e := io.Copy(this.file, &this.buffer)
		this.curSize += n
		if e != nil {
			Print(os.Stderr, LevelError, 0, FileLineFullPath, e.Error())
		}
		// 写入的数据到达最大，开始新文件
		if this.curSize >= int64(this.maxSize) {
			this.newFile()
		}
	}
	this.mutex.Unlock()
}

// 将缓存的数据保存到磁盘文件中
func (this *LoggerFile) newFile() {
	// 关闭旧文件
	this.closeFile()
	// 时间
	now := time.Now()
	// 创建日期目录
	date_dir := filepath.Join(this.rootDir, now.Format(this.dayFormat))
	e := os.MkdirAll(date_dir, os.ModePerm)
	if nil != e {
		Print(os.Stderr, LevelError, 1, FileLineFullPath, e.Error())
		return
	}
	// 新的日志文件
	time_file := filepath.Join(date_dir, now.Format(this.fileFormat))
	this.file, e = os.OpenFile(time_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != e {
		Print(os.Stderr, LevelError, 1, FileLineFullPath, e.Error())
		return
	}
}

// 关闭文件
func (this *LoggerFile) closeFile() {
	if nil != this.file {
		// 把剩下的缓存写入文件，再关闭
		io.Copy(this.file, &this.buffer)
		this.file.Close()
		this.file = nil
	}
}

func (this *LoggerFile) Close() error {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return errors.New("logger has been closed")
	}
	this.closed = true
	this.mutex.Unlock()
	// 退出
	this.syncTimer.Reset(0)
	// 等待routine退出
	<-this.exit
	// 关闭文件
	this.closeFile()
	return nil
}

// 新的日志文件对象
func NewFileLogger(cfg *LoggerFileConfig) *LoggerFile {
	// 解析配置参数
	size, e := common.ParseInt(cfg.Size)
	if nil != e {
		size = 1024 * 1024
	}
	dur, e := time.ParseDuration(cfg.Duration)
	if nil != e {
		dur = time.Second * 3
	}
	// 对象
	lf := &LoggerFile{
		rootDir:    cfg.Dir,
		exit:       make(chan struct{}),
		maxDay:     common.MaxInt(cfg.Day, 1),
		maxSize:    common.MaxInt(int(size), 1024*1024),
		duration:   common.MaxDuration(dur, time.Second),
		dayFormat:  cfg.DayFormat,
		fileFormat: cfg.FileFormat,
	}
	lf.SetStdWriter(cfg.Std)
	if lf.rootDir == "" {
		lf.rootDir = "./"
	}
	if lf.dayFormat == "" {
		lf.dayFormat = "20060102"
	}
	if lf.fileFormat == "" {
		lf.fileFormat = "150405.999999999"
	}
	lf.newFile()
	// 保存routine
	go func(this *LoggerFile) {
		defer Recover(this, true, false, func() {
			this.syncTimer.Stop()
			close(this.exit)
		})
		this.syncTimer = time.NewTimer(this.duration)
		for !this.closed {
			<-this.syncTimer.C
			// 保存数据到文件
			this.saveFile()
			// 检查过期的日志
			this.checkExpired()
			this.syncTimer.Reset(this.duration)
		}
	}(lf)

	return lf
}
