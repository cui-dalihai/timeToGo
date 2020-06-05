package main

import (
	"fmt"
)

type Func func(key string) (interface{}, error)
type result struct {
	value interface{}
	err error
}
type entry struct {
	res result
	ready chan struct{}
}
func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}
func (e *entry) deliver(response chan <- result) {
	<- e.ready
	response <- e.res
}
//type Memo struct {
//	f Func
//	mu sync.Mutex
//	cache map[string]*entry
//}
//func New(f Func) *Memo {
//	return &Memo{f: f, cache:make(map[string]*entry)}
//}
//func (memo *Memo) Get(key string) (value interface{}, err error) {
//	memo.mu.Lock()
//	e := memo.cache[key]
//	if e == nil {
//		e = &entry{ready: make(chan struct{})}
//		memo.cache[key] = e
//		memo.mu.Unlock()
//		e.res.value, e.res.err = memo.f(key)
//		close(e.ready)
//	} else {
//		memo.mu.Unlock()
//		<- e.ready
//	}
//	return e.res.value, e.res.err
//}

type request struct {
	key string
	response chan <- result
}
type Memo struct {
	requests chan request
}
func (memo *Memo) Close() { close(memo.requests) }
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <- response
	return res.value, res.err
}
func (memo *Memo) Server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready:make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}
func New(f Func) *Memo {
	memo := &Memo{requests:make(chan request)}
	go memo.Server(f)
	return memo
}




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
