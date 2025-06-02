package main

import (
	"bufio"
	"fmt"
	"net"
)

type Message struct {
	Connection net.Conn
	Content    string
	UserName   string
}

var clients = make(map[net.Conn]string)
var messages = make(chan Message)

func handleConnection(conn net.Conn) {

	defer conn.Close()

	clients[conn] = conn.RemoteAddr().String()

	name := conn.RemoteAddr().String()

	fmt.Printf("New connection from %s\n", name)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("%s # %s\n", name, message)

		messages <- Message{Connection: conn, Content: message, UserName: name}

	}

	delete(clients, conn)
}

func broadcast() {
	for msg := range messages {
		for client := range clients {
			if client == msg.Connection {
				continue
			}

			fmt.Fprintf(client, "%s @ %s\n", msg.UserName, msg.Content)
		}
	}
}

func main() {

	listener, err := net.Listen("tcp", ":45956")
	if err != nil {
		fmt.Printf("Failed to start listener: %s\n", err)
	}
	defer listener.Close()

	go broadcast()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}

}
