package main

import (
	"bytes"
	"fmt"
)

var months = [...]string{1:"Jan", 2:"Feb", 3:"Mar", 4:"Apr", 5:"May", 6:"Jun", 7:"Jul", 8:"Aug", 9:"Sep", 10:"Oct", 11:"Nov", 12:"Dec"}

func main() {

	Q2 := months[4:7]
	summer := months[6:9]

	fmt.Println(Q2)
	fmt.Println(summer)
	fmt.Println(len(summer))
	es := summer[:5]
	fmt.Println(es)
	fmt.Println(len(summer))






	//var s = "hello, 世界"
	//var t = "\u4e16\uffff"
	//var m = '\ufff2'
	//fmt.Println(len(s))  // the number of bytes in s
	//fmt.Println(utf8.RuneCountInString(s))
	//fmt.Printf("%q\n", t)
	//fmt.Printf("%q\n", m)
	//fmt.Println(reflect.TypeOf(s))
	//fmt.Println(reflect.TypeOf(t))
	//fmt.Println(reflect.TypeOf(m))
	//
	//var f = "test abc 和我"
	//c := []byte(f)
	//b := string(c)
	//fmt.Println(f)
	//fmt.Println(reflect.TypeOf(f))
	//fmt.Println(c)
	//fmt.Println(reflect.TypeOf(c))
	//fmt.Println(b)
	//fmt.Println(reflect.TypeOf(b))
	//
	//for i, r := range s {
	//	fmt.Printf("%d\t%c\t%d\n", i, r, r)
	//}
	//
	//fmt.Println(intsToStr([]int{1,2,3,4,5}))
	//fmt.Println(comma1("81231231"))
}

func reverse(s []int) {
	for i, j := 0, len(s) - 1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func intsToStr (values []int) string {
	var buf bytes.Buffer
	buf.WriteRune('[')

	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteRune(']')
	return buf.String()
}

func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func comma1(s string) string {
	n := len(s)
	pr := n % 3

	var buf bytes.Buffer

	for i, v := range s {
		ia := i % 3
		if i == n {
			break
		}
		if ia == pr && (i != 0) {
			buf.WriteRune(',')
		}
		buf.WriteRune(v)
	}

	return buf.String()
}

func zeroArray(ptr *[32]byte) {
	for i := range ptr {
		ptr[i] = 0
	}
}
func zeroA(ptr *[32]byte) {
	*ptr = [32]byte{}
}

