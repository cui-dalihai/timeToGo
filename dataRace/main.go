package main

import (
	"fmt"
	"image"
	"sync"
)

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) {deposits <- amount}
func Balance() int {return <- balances}

func teller() {
	var balance int
	for {
		select {
		case amount := <- deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}

var (
	sema = make(chan struct{}, 1)  // 容量为1
	balance1 int
)

func Deposit1(amount int) {
	sema <- struct{}{}  // 写入成功表示获取了权限
	balance1 += amount
	<- sema
}

func Balance1() int {
	sema <- struct{}{}
	b := balance1
	<- sema
	return b
}

var (
	mu  sync.Mutex
	balance2 int
)

var (
	mu3 sync.RWMutex
	balance3 int
)

func Balance3() int {
	mu3.RLock()
	defer mu3.RUnlock()
	return balance3
}

func Deposit2(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func deposit(amount int) {balance2 += amount}

func Balance2() int {
	mu.Lock()
	defer mu.Unlock()
	return balance2
}


func Withdraw2(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance2 < 0 {
		deposit(amount)
		return false
	}
	return true
}


var icons map[string]image.Image

func loadIcon(icon_name string) image.Image {
	return nil
}

func loadIcons() {
	icons = map[string]image.Image{
		"spades.png": loadIcon("spades.png"),
		"hearts.png": loadIcon("hearts.png"),
	}
}

func Icon(name string) image.Image {
	if icons == nil {
		loadIcons()
	}
	return icons[name]
}

var mu1 sync.RWMutex

func Icon1(name string) image.Image {
	mu1.RLock()
	if icons != nil {
		icon := icons[name]
		mu1.RUnlock()
		return icon
	}
	mu1.RUnlock()   // 处理else的情况

	mu1.Lock()
	if icons == nil {
		loadIcons()
	}
	icon := icons[name]
	mu1.Unlock()
	return icon
}

var loadIconsOnce sync.Once

func Icon2(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}


func main() {
	Deposit2(100)
	Deposit2(200)
	Deposit2(300)
	fmt.Println(Balance2())
	fmt.Println(Withdraw2(700))
	fmt.Println(Balance2())
	fmt.Println(Withdraw2(200))
	fmt.Println(Balance2())
}

