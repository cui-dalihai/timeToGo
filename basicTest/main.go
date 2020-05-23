package main

import (
	"fmt"
)

func incomingURLs () []string {
	l := []string{
		"https://golang.org",
		"https://godoc.org",
		"https://play.golang.org",
		"http://gopl.io",
		"https://golang.org",
		"https://godoc.org",
		"https://play.golang.org",
		"http://gopl.io",
	}
	return l
}

func main()  {
	allUrls := incomingURLs()
	for index := range allUrls {
		fmt.Printf("%s, \n", allUrls[index])
	}
}
