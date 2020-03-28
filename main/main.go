package main

import (
	"go-logger/mylogger"
	"time"
)

func main() {
	// log := mylogger.NewLog("info")
	log := mylogger.NewFileLogger("info", "./", "test.log", 10*1024*1024)
	for {
		log.Debug("这是一条debug日志")
		id := 100
		name := "ccc"
		log.Info("这是一条info日志,id:%v name:%s", id, name)
		log.Warning("这是一条warning日志")
		log.Error("这是一条error日志")
		log.Fatal("这是一条fatal日志")
		time.Sleep(time.Second)
	}
}
