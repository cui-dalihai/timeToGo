package memo_test

import (
	"concurentNonBlockingCache/memo"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
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
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

type M interface {
	Get(key string) (interface{}, error)
}

//func Seqs(t *testing.T, m M){
//	for url := range incomingURLs() {
//		start := time.Now()
//		value, err := m.Get1(url)
//		if err != nil {
//			log.Print(err)
//			continue
//		}
//		fmt.Printf("%s, %s, %d bytest\n", url, time.Since(start), len(value.([]byte)))
//	}
//}

func Conc(t *testing.T, m M) {
	var n sync.WaitGroup
	for url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytest\n", url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	n.Wait()
}


//func Test(t *testing.T) {
//	m := memo.New(httpGetBody)
//	Seqs(t, m)
//}

// NOTE: not concurrency-safe!  Test fails.
func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	Conc(t, m)
}
