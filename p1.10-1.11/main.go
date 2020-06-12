package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	// 阻塞等待最后一个响应返回
	for range os.Args[1:] {
		fmt.Println(<- ch)
	}

	// 取决于最长时间的那个请求
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	// 直接丢弃结果, 只统计字节数
	//nb, err := io.Copy(ioutil.Discard, resp.Body)
	f, err := os.Create("test.html")
	if err != nil {
		ch <- fmt.Sprintf("while creating %s: %v\n", url, err)
		return
	}

	nb, err := io.Copy(f, resp.Body)

	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs	%7d	%s\n", secs, nb, url)
}

