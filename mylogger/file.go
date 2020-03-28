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

func (f *FileLogger) enable(LogLevel LogLevel) bool {
	return LogLevel >= f.Level
}

func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n",
			now.Format("2006-01-02 15:04:05"), getLogstring(lv),
			fileName, funcName, lineNo, msg)
		if lv >= ERROR {
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
