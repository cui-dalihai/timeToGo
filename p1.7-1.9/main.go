package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {

		if strings.HasPrefix(url, "http://") {
			fmt.Printf("Already satisfied.\n")
		} else {
			fmt.Printf("Prefixing...\n")
			url = "http://" + url
			fmt.Printf("Prefixed: %s\n", url)
		}


		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "1.7fetch: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("HTTP Status Code: %d\n", resp.StatusCode)

		//b, err := ioutil.ReadAll(resp.Body)

		// copy函数尝试src.WriteTo(dst)和dst.ReadFrom(src)两种方式
		b, err := io.Copy(os.Stdout, resp.Body)

		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "1.7fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%d\n", b)
	}
}

