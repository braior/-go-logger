package mylogger

import (
	"fmt"
	"time"
)

// 往终端上写日志相关内容

// ConsoleLogger 日志结构体
type ConsoleLogger struct {
	Level LogLevel
}

// NewLog 日志结构体构造函数
func NewConsoleLogger(leveStr string) ConsoleLogger {
	level, err := parseLogLevel(leveStr)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		Level: level,
	}
}

func (c ConsoleLogger) enable(LogLevel LogLevel) bool {
	return LogLevel >= c.Level
}

func (c ConsoleLogger) log(lv LogLevel, format string, a ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n",
			now.Format("2006-01-02 15:04:05"), getLogstring(lv),
			fileName, funcName, lineNo, msg)
	}
}

// Debug 记录debug日志函数
func (c ConsoleLogger) Debug(format string, a ...interface{}) {

	c.log(DEBUG, format, a...)

}

// Info 记录info日志函数
func (c ConsoleLogger) Info(format string, a ...interface{}) {
	c.log(INFO, format, a...)

}

// Warning 记录Warning日志函数
func (c ConsoleLogger) Warning(format string, a ...interface{}) {
	c.log(WARNING, format, a...)
}

// Error 记录Error日志函数
func (c ConsoleLogger) Error(format string, a ...interface{}) {
	c.log(ERROR, format, a...)
}

// Fatal 记录Fatal日志函数
func (c ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.log(FATAL, format, a...)
}
