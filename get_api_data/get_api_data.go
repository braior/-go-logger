package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Data ..
type Data struct {
	TradeID string `json:"tradeID"`
	Date    string `json:"date"`
	Type    string `json:"type"`
	Rate    string `json:"rate"`
	Amount  string `json:"amount"`
}

// Result ..
type Result struct {
	Res     string `json:"result"`
	Da      []Data `json:"data"`
	Elapsed string `json:"elapsed"`
}

//请求API获取返回值，并把json解析为struct
func getAPI(pair string) (ResultPTR *Result, err error) {
	var apiURL = "https://api.ukex.io/api/index/tradeHistory/" + pair

	var ResultByte []byte
	response, err := http.Get(apiURL)
	if err != nil {
		return ResultPTR, err
	}

	ResultByte, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return ResultPTR, err

	}

	ResultPTR = &Result{}
	err = json.Unmarshal(ResultByte, ResultPTR)
	return ResultPTR, err
}

//字符串类型时间转换为时间类型
func timeTrans(date string) time.Time {
	date800 := date + " +0000"
	Time, err := time.Parse("2006-01-02 15:04:05 -0700", date800)
	if err != nil {
		fmt.Println(1202)
	}
	return Time
}

//获取系统现在的时间跟API请求的最新时间比较
func timeCompare(tradeTime time.Time) int64 {
	localTime := time.Now().Unix()
	apiTime := tradeTime.Unix()

	if apiTime <= localTime {
		return localTime - apiTime
	}
	return 1201
}

func main() {
	var pair string
	if len(os.Args) == 2 {
		pair = os.Args[1]
		respon, err := getAPI(pair)
		if err != nil {
			fmt.Println(1203)
			fmt.Println(err)
		} else {
			fmt.Println(respon)
		}
		tradeTime := timeTrans(respon.Da[0].Date)
		fmt.Println(timeCompare(tradeTime))
	} else {
		fmt.Println(1204)
	}
}

