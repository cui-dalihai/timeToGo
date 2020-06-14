package main

import (
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

func bInt8(n int8) string {
	return strconv.FormatUint(uint64(*(*uint8)(unsafe.Pointer(&n))), 2)
}
func bInt64(n int64) string {
	return strconv.FormatUint(uint64(*(*uint64)(unsafe.Pointer(&n))), 2)
}
//func bFloat32(n float32) string {
//	return strconv.FormatUint(uint32(*(*uint32)(unsafe.Pointer(&n))), 2)
//}

func main() {
	//var s uint
	//var b int
	//var t int64
	//var f int16
	//f = 2147483646
	//f = math.MaxInt16 - 1
	//f += 1
	//
	//h := uint8(f)

	//var g int8   // -128 ~ 127
	//g = -128
	//fmt.Printf("%d\n", g)
	//fmt.Printf("%b\n", g)
	//g -= 1
	//fmt.Printf("%d\n", g)
	//fmt.Printf("%b\n", g)


	//var i uint8  // 0 ~ 255
	//i = 255
	//i += 1  // 加1溢出
	//i = 0
	//i -= 1  // 减1溢出


	//fmt.Printf("%d\n", unsafe.Sizeof(s))
	//fmt.Printf("%d\n", unsafe.Sizeof(b))
	//fmt.Printf("%d\n", unsafe.Sizeof(t))
	//fmt.Printf("%d\n", unsafe.Sizeof(f))
	//fmt.Printf("%b\n", f)
	//fmt.Printf("%d\n", f)
	//fmt.Printf("%b\n", h)
	//fmt.Printf("%d\n", h)
	//fmt.Printf("%d\n", unsafe.Sizeof(g))
	//fmt.Printf("%d\n", unsafe.Sizeof(i))
	//fmt.Printf("%b\n", i)

	//var al byte
	//var a1 = []byte{1}
	//fmt.Printf("%T\n", a1)
	//fmt.Printf("%d\n", al)
	//fmt.Println(a1[0])

	var a, b rune
	a = math.MaxInt32 - 5
	b = math.MinInt32

	c := int8(a)                    // 低位截断: 11111010  = -1 * 2^7 + 2^6 + 2^5 + 2^4 + 2^3 + 2 = -6
	fmt.Printf("%d\n", c)
	fmt.Printf("低位截断: %s\n", bInt8(c))

	e := c << 2                     // 左移位: 11101000 即 -6 * 2^2 = -24
	fmt.Printf("%d\n", e)
	fmt.Printf("左移: %s\n", bInt8(e))

	d := int64(b)                   // 负数高位补1扩展
	fmt.Printf("%d\n", d)
	fmt.Printf("负数高位补1扩展: %s\n", bInt64(d))

	f := c + 8                      // 11111010 + 1000 = 100000010 截断 00000010 = 2
	fmt.Printf("%d\n", f)
	fmt.Printf("加法截断: %s\n", bInt8(f))

	fo := c + (-125)               // 1111 1010 + 1000 0011 = 1 0111 1101 截断 0111 1101 即 125
	fmt.Printf("%d\n", fo)
	fmt.Printf("加法溢出截断: %s\n", bInt8(fo))

	g := c * 5                      // 11111010 * (2^2 + 1) : 左移两位再加上之前的自己 1110 1000 + 1111 1010 = 1 1110 0010 截断后为 1110 0010 即-30
	fmt.Printf("%d\n", g)
	fmt.Printf("乘法截断: %s\n", bInt8(g))

	gof := c * 100                  // 11111010 * (2^6+2^5+2^2): 1000 0000 + 0100 0000 + 1110 1000 = 1 1010 1000 截断后为 1010 1000 即 -88
	fmt.Printf("%d\n", gof)
	fmt.Printf("乘法溢出截断: %s\n", bInt8(gof))

	m := c / -2  // 6/2	            // (0000 0110)右移一位: 0000 0011 即 3
	fmt.Printf("%d\n", m)
	fmt.Printf("除法: %s\n", bInt8(m))

	j := c / 4   // c为负数,需要偏置  // (11111010 + 4 - 1) * (2^-2): (11111010 + 0000 0011)= 1111 1101 再右移两位 1111 1111 即 -1
	fmt.Printf("%d\n", j)
	fmt.Printf("除法偏置: %s\n", bInt8(j))

	//mf := float32(c)
	var mf float32 = -24.125                                    // 1 10000011 10000010000000000000000: 符号s=1位, 阶码k=8位， 尾数n=23位,
	fmt.Printf("%f\n", mf)                                      // 阶码既非全0也非全1, 规格化类型, Bias = 2^(k-1) - 1 = 2^7 -1
                                                                // e = U(10000011) = 2^7 + 2^1 + 2^0 = 2^7 + 3
                                                                // E = e - Bias = 2^7 + 3 - 2^7 + 1 = 4
	fmt.Printf("规格化的浮点数: %b\n", math.Float32bits(mf))      // M = 1 + f = 1 + 2^-1 + 2^-7
                                                                // V = (-1)^s * M * 2^E = -1 * (1 + 2^-1 + 2^-7) * 2^4 = -(24 + 1/8)

                                                                // 1 00000000 11111111111111111111111
                                                                // 非规格化时, 阶码全为0: M = f; E = 1 - Bias
                                                                // M = f = 2^-1 + 2^-2 + 2^-3 + ... + 2^-23
                                                                // E = 1 - Bias = 1 - 2^7 + 1 = 2 - 2^7
                                                                // V = -1 * (2^-1 + 2^-2 + 2^-3 + ... + 2^-23) * 2^(2 - 2^7) ≈ 2^-126 ≈ 1.1754942E-38
	var imf float32 = -1.1754942E-38                            // 实际上, 去掉负号这是32位float所能表示的最大的非规格数, 因为再大一些所能表示的浮点数是
	fmt.Printf("%f\n", imf)                                     // 0 00000001 000000000000000000000000, 这将进入规格化的算法范围
	fmt.Printf("非规格化浮点数: %b\n", math.Float32bits(imf))
                                                                // 特殊值: 阶码全为1时, 小数全为0为无穷, s为1和0分别表示-∞和+无穷
                                                                //                   小数是非零, 表示NaN: not a number

}