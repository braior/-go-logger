package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	//\d代表数字
	reQQEmail = `(\d+)@qq.com`

	//匹配邮箱,()内的匹配项为一个子表达式
	//reMail = `\w+@\w+\.\w+(\.\w+)?`
	reMail = `\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}`

	//链接，(.+?)为惰性匹配
	reLink = `href="(https?://[\s\S]+?)"`

	//手机号码
	//rePhone = `1[3456789]\d\s?\d{4}\s?\d{4}`
	rePhone = `(13\d|14[579]|15[^4\D]|17[^49\D]|18\d)\d{8}`
	//410222 1987 06 13 4038
	reIdcard = `[12345678]\d{5}((19\d{2})|(20[01]))((0[1-9]|[1[012]]))((0[1-9])|[12]\d|[3[01]])\d{3}[\dXx]`
	//reIdcard = ` ^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
	//图片链接
	reImg = `"(https?://[^"]+?(\.(jpg|jpeg|png|gif|ico)))"`
	//reImg = `https?:\/\/(.+\/)+.+(\.(gif|png|jpg|jpeg|webp|svg|psd|bmp|tif))$`
)

//HandleError 处理异常
func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

//GetPageStr 根据url获取页面内容
func GetPageStr(url string) (pageStr string) {
	//发送http请求，获取页面内容
	resp, err := http.Get(url)
	//处理异常
	HandleError(err, "http.Get url")
	//关闭资源
	defer resp.Body.Close()
	//接收页面内容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	//打印页面内容
	pageStr = string(pageBytes)
	// fmt.Println(pageStr)
	return pageStr

}

// GetEmail 爬取邮箱
func GetEmail(url string) {
	//调用GetPageStr函数， 根据url获取页面内容
	pageStr := GetPageStr(url)

	//捕获邮箱，先搞定qq邮箱
	//传入正则表达式，得到正则表达式对象
	re := regexp.MustCompile(reQQEmail)
	results := re.FindAllStringSubmatch(pageStr, -1)
	for _, result := range results {
		fmt.Printf("email=%s qq=%s\n", result[0], result[1])
	}

}

//GetEmail2 抽取的爬邮箱的方法
func GetEmail2(url string) {
	//爬取页面所有数据
	pageStr := GetPageStr(url)

	//传入正则表达式，得到正则表达式对象
	re := regexp.MustCompile(reMail)
	fmt.Println(re)

	// 用正则对象提取页面中内容，pageStr是页面内容，-1代表取所有
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Println(results)
	for _, result := range results {
		fmt.Println(result[0][0])
	}
}

//GetLink 爬超链接
func GetLink(url string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reLink)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("找到%d条结果:\n", len(results))
	for _, result := range results {
		//fmt.Println(result)
		fmt.Println(result[0])
	}
}

//GetPhone 爬取手机号
func GetPhone(url string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(rePhone)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("找到%d条结果:\n", len(results))
	for _, result := range results {
		fmt.Println(result[0])
	}
}

//GetIdcard 爬取身份证
func GetIdcard(url string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reIdcard)

	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("找到%d条结果:\n", len(results))
	fmt.Println(results)
	for _, result := range results {
		fmt.Println(result[0])
	}
}

//GetImgURL 获取图片链接
func GetImgURL(url string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("找到%d条结果:\n", len(results))
	for _, result := range results {
		fmt.Println(result[0])
	}
}

var (
	// 存放图片链接
	chanImgUrls chan string

	// 存放147个任务，判断任务是否已完成
	chanTask  chan string
	waitGroup sync.WaitGroup
)

// SpiderPrettyImg 爬取一个页面上的全部图片链接，返回结果切片
func SpiderPrettyImg(url string) (urls []string) {
	pageStr := GetPageStr(url)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Println(results)
	for _, result := range results {
		url := result[1]
		urls = append(urls, url)
	}
	//fmt.Println(urls)
	return
}

//PutImgURLToChan 将爬取的链接加到管道
func PutImgURLToChan(url string) {
	//爬取当前页面所有图片链接
	urls := SpiderPrettyImg(url)
	fmt.Println(urls)

	//添加到管道
	for _, url := range urls {
		chanImgUrls <- url
	}

	// 关闭通道
	close(chanImgUrls)

	//标记当前携程任务已完成
	//chanTask <- url
	//waitGroup.Done()
}

// GetFilenameFromURL 从url中提取内容并拼接文件名
func GetFilenameFromURL(url string, dirPath string) (filename string) {
	//strings包的方法，截取最后一个/
	lastIndex := strings.LastIndex(url, "/")
	filename = url[lastIndex+1:]
	//加一个时间戳，防止重名
	// timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	// filename = timePrefix + "_" + filename
	filename = dirPath + filename
	return
}

//DownloadFile 下载file对应的文件到指定路径
func DownloadFile(url string, filename string) (ok bool) {

	resp, err := http.Get(url)
	if err != nil {
		HandleError(err, "http.Get(url)")
		return
	}
	defer resp.Body.Close()

	//ioutil.ReadAll(resp.Body)read tcp 192.168.20.50:57178->175.6.244.4:80: wsarecv:
	// An existing connection was forcibly closed by the remote host.

	fBytes, e := ioutil.ReadAll(resp.Body)
	HandleError(e, "ioutil.ReadAll(resp.Body)")
	err = ioutil.WriteFile(filename, fBytes, 0644)
	HandleError(err, "ioutil.WriteFile(filename, fBytes, 0644)")
	if err != nil {
		return false
	}
	return true
}

// DownloadImg 同步下载图片链接管道中的所有图片
func DownloadImg(i int) {
	for url := range chanImgUrls {
		filename := GetFilenameFromURL(url, "D:/gopath/src/go_project/concurrent_reptiles/images/")
		ok := DownloadFile(url, filename)

		if ok {
			fmt.Printf("进程%d开始下载，%s下载成功！\n", i, filename)
		} else {
			fmt.Printf("进程%d开始下载，%s下载失败！\n", i, filename)
		}
	}
	waitGroup.Done()
}

// CheckIfAllSpiderOK 检查147个任务是否全部完成，完成则关闭数据管道
func CheckIfAllSpiderOK() {
	var count int
	for {
		url := <-chanTask
		fmt.Printf("%s完成爬取任务\n", url)
		count++
		if count == 147 {
			close(chanImgUrls)
			break
		}
	}
	waitGroup.Done()
}

func main() {

	//var url = "https://tieba.baidu.com/p/2366935784"

	//var phoneURL = "https://www.zhaohaowang.com/"
	//1.爬取邮箱
	//GetEmail(url)
	//2.抽取爬邮箱的方法
	//GetEmail2(url)
	//3.爬超链接
	//GetLink(url)
	//4.爬手机号
	//GetPhone(phoneURL)
	//5.爬身份证
	//GetIdcard("http://henan.qq.com/a/20171107/069413.htm")
	//6.爬图片链接
	//GetImgURL("http://image.baidu.com/search/index?tn=baiduimage&ps=1&ct=201326592&lm=-1&cl=2&nc=1&ie=utf-8&word=%E7%BE%8E%E5%A5%B3")

	//初始化数据管道
	chanImgUrls = make(chan string, 1000000)
	chanTask = make(chan string, 147)

	//爬虫协程：远远不断的往管道中添加图片链接
	// for i := 1; i < 148; i++ {
	// 	waitGroup.Add(1)
	// 	//获取某个页面所有图片链接
	// 	//strconv.Itoa(i)：将整数转为字符串
	// 	go PutImgURLToChan("https://www.enterdesk.com/special/wmtp/fzlwm/")
	// }

	//爬取单个页面所有图片链接并发送到通道中
	PutImgURLToChan("https://www.enterdesk.com/special/wmtp/fzlwm/")

	// 开辟任务统计协程，如果147个任务全部完成，则关闭数据管道
	// waitGroup.Add(1)
	// go CheckIfAllSpiderOK()

	//下载协程：源源不断从管道中读取地址并下载
	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go DownloadImg(i)
	}
	waitGroup.Wait()
}
