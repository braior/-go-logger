package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里面写日志

// FileLogger struct
type FileLogger struct {
	Level       LogLevel
	filePtah    string // 日志文件保存的路径
	fileName    string // 日志文件名
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
}

// NewFileLogger 构造函数
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	LogLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	f1 := &FileLogger{
		Level:       LogLevel,
		filePtah:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}
	err = f1.initFile() //按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return f1
}

// 根据指定的日志路径和文件名打开日志文件
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePtah, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open err log file failed,err:%v\n", err)
		return err
	}
	// 日志文件都已经打开了
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

// 判断是否需要记录该日志
func (f *FileLogger) enable(LogLevel LogLevel) bool {
	return LogLevel >= f.Level
}

// 判断文件大小是否切割
func (f *FileLogger) checkLogFileSize(fileLog *os.File) bool {
	fileInfo, err := fileLog.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err:%v\n", err)
		return false
	}
	// 判断文件大小是否大于设定的最大值
	return fileInfo.Size() >= f.maxFileSize

}

// 切割日志文件
func (f *FileLogger) splitLogFile(file *os.File) (*os.File, error) {
	// 需要切割日志文件
	nowStr := time.Now().Format("2006010215040500")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err:%v", err)
		return nil, err
	}
	logName := path.Join(f.filePtah, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.%s", logName, nowStr)
	// 1.关闭日志文件
	file.Close()
	// 2.备份一下 rename
	os.Rename(logName, newLogName)
	// 3.打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed, err:%v", err)
		return nil, err
	}
	// 4.打开新的日志文件对象赋值给 f.fileObj
	return fileObj, nil
}

func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		if f.checkLogFileSize(f.fileObj) {
			newLogFile, err := f.splitLogFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newLogFile
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n",
			now.Format("2006-01-02 15:04:05"), getLogstring(lv),
			fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			if f.checkLogFileSize(f.errFileObj) {
				newLogFile, err := f.splitLogFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newLogFile
			}
			// 如果要记录的日志大于等于ERROR级别，将在err日志文件再记录一遍
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n",
				now.Format("2006-01-02 15:04:05"), getLogstring(lv),
				fileName, funcName, lineNo, msg)
		}
	}
}

// Debug 记录debug日志函数
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

// Info 记录info日志函数
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

// Warning 记录Warning日志函数
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

// Error 记录Error日志函数
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

// Fatal 记录Fatal日志函数
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}

// Close file
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
