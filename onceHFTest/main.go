package main

import "sync"

var once sync.Once
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func doprint() {
	print("done:", done)
	if !done {
		once.Do(setup)
	}
	print("this:", a)
}

func main() {
	go doprint()
	go doprint()
}

