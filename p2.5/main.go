package main

import "fmt"

func main() {
	x := 2
	fmt.Println(x)
	fmt.Println(&x)
	s := &x
	fmt.Println(s)
	fmt.Println(&s)
	t := &s
	fmt.Println(t)
	fmt.Println(&t)
	fmt.Println(**t)
}

