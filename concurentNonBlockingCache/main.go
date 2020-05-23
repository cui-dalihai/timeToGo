package main

import (
	"concurentNonBlockingCache/memo"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func incomingURLs () []string {
	l := []string{
		"https://baidu.com",
		"https://www.cctv.com/",
		"https://sina.com.cn",
		"http://www.bing.com/?mkt=zh-CN",
		//"https://weibo.com",
		//"https://www.csdn.net",
		//"https://xueqiu.com",
		//"https://www.jianshu.com/",
		//"https://www.zhihu.com/",
		"https://baidu.com",
		"https://www.cctv.com/",
		"https://sina.com.cn",
		//"http://www.bing.com/?mkt=zh-CN",
		//"https://weibo.com",
		//"https://www.csdn.net",
		//"https://xueqiu.com",
		//"https://www.jianshu.com/",
		//"https://www.zhihu.com/",
	}
	return l
}

func HttpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main()  {

	// 串行执行
	//m := memo.New(httpGetBody)
	//allUrls := incomingURLs()
	//for i := range allUrls {
	//	url := allUrls[i]
	//	start := time.Now()
	//	value, err := m.Get1(url)
	//	if err != nil {
	//		log.Print(err)
	//	}
	//	fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
	//}


	// 并发执行, 为每个url单独创建一个goroutine来请求
	// 由于get请求是耗时操作, 那么在缓存建立之前, 所有的goroutine可能都已经完成缓存检查且结果是都没有缓存, 所有的goroutine都去get
	// 并且相同url的两个goroutine后响应的结果会覆盖前者, -race 会检测出这两个goroutine在未同步的情况下写了相同位置, 即出现了竞态
	m := memo.New(HttpGetBody)
	var n sync.WaitGroup
	allUrls := incomingURLs()
	for i := range allUrls {
		url := allUrls[i]
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
			n.Done()
		}(url)
	}
	// block当前goroutine,直到WaitGroup中为0
	n.Wait()
}

