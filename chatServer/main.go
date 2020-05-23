package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleConn(conn)
	}
}

// client为一个单向通道类型，只允许向其中发送string类型数据
type client chan <- string
var (
	entering = make(chan client)  // client类型的双向通道
	leaving = make(chan client)
	messages = make(chan string)
)

func broadcaster () {
	// 使用短变量声明来声明和初始化局部变量, key是client类型, value是布尔类型
	clients := make(map[client]bool)
	for {
		select {
		case msg := <- messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <- entering:
			print("someone has arrived.\n")
			clients[cli] = true

		case cli := <- leaving:
			print("someone has left.\n")
			delete(clients, cli)
			close(cli)
		}

	}
}

func handleConn(conn net.Conn)  {
	ch := make(chan string)

	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived."
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who	+ ":" + input.Text()
	}

	leaving <- ch
	messages <- who + " has left."
	err := conn.Close()
	if err != nil {
		log.Print(err)
	}
}

func clientWriter(conn net.Conn, ch <- chan string)  {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}