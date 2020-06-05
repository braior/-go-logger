package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//分析nginx log日志，截取日志的各个字段，并统计各个IP的访问量并倒叙排序，计算出当前日志的QPS

// LogMap ..
var LogMap map[string]int

// SortLogMapOrderbyNum ..
var SortLogMapOrderbyNum map[string]int

// IPNum 结构体 记录IP地址和其出现的次数
type IPNum struct {
	IPAddress string
	Num       int
}

// timeTihuan 将日志里的时间转换为时间戳
func timeTihuan(dateHour string) int64 {
	theTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", dateHour)
	if err != nil {
		fmt.Println(err)
	}
	return theTime.Unix()
}

func readLog() int64 {

	LogMap = make(map[string]int)
	var firstTime, lastTime, durTime int64

	file, err := os.Open("1.txt")
	defer file.Close()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return 1
	}

	// 读取文件内容到bufio
	line := bufio.NewReader(file)
	var flag bool = true
	for {

		// 逐行读取
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}

		contentStr := string(content)
		lineSplit := strings.Split(contentStr, " ")

		_, ok := LogMap[lineSplit[0]]
		if ok {
			LogMap[lineSplit[0]]++
		} else {
			LogMap[lineSplit[0]] = 1
		}

		//获取日志时间
		logTimeString := lineSplit[3] + " " + lineSplit[4]
		logTimeUnix := timeTihuan(logTimeString)
		//logTimeStringUnix := strconv.FormatInt(logTimeUnix, 10)
		// fmt.Println(logTimeStringUnix)

		if flag {
			firstTime = logTimeUnix
			flag = false
		} else {
			lastTime = logTimeUnix
		}
	}
	durTime = lastTime - firstTime
	// fmt.Println(durTime)

	fmt.Println(LogMap)
	return durTime
}

func sortIPOrderByNum() {

	var lstIPNum []IPNum
	for k, v := range LogMap {
		lstIPNum = append(lstIPNum, IPNum{k, v})
	}

	sort.Slice(lstIPNum, func(i, j int) bool {
		return lstIPNum[i].Num > lstIPNum[j].Num
	})

	// fmt.Println(lstIPNum)

	for i := range lstIPNum {
		fmt.Printf("ip: %s \t number: %d\n", lstIPNum[i].IPAddress, lstIPNum[i].Num)
	}
}

// Decimal 将float64数据类型保留小数位
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func main() {
	var sum int
	durTime := readLog()
	fmt.Println(durTime)
	sortIPOrderByNum()

	for k := range LogMap {
		sum += LogMap[k]
	}
	fmt.Println(sum)

	QPS := float64(sum) / float64(durTime)
	fmt.Println(QPS)
}
