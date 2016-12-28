package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

type clientInfo struct {
	nick      string
	logonTime time.Time
}

type message struct {
	c       net.Conn
	payload string
}

func main() {
	fmt.Println("Starting server on port 6000")

	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conns := make(chan net.Conn, 10)
	go dispatcher(conns)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		conns <- conn
	}
}

func dispatcher(connections <-chan net.Conn) {
	clients := make(map[net.Conn]clientInfo)
	msgPipe := make(chan message, 5)
	i := 1
	for {
		select {
		case conn := <-connections:
			clients[conn] = clientInfo{nick: "client" + strconv.Itoa(i), logonTime: time.Now()}
			i++
			go handleConnection(conn, msgPipe)
		case msg := <-msgPipe:
			fmt.Println(msg)
			info := clients[msg.c]
			for c := range clients {
				if c != msg.c {
					c.Write([]byte(info.nick + "> " + msg.payload + "\n"))
				}
			}
		}
	}
}

func handleConnection(conn net.Conn, pipe chan<- message) {
	var msg message
	msg.c = conn
	fmt.Printf("Got connection from %s\n", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg.payload = scanner.Text()
		pipe <- msg
	}
	fmt.Println("Client disconnected.")
	conn.Close()
}
