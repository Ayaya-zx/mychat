package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type client struct {
	ch   chan string
	conn net.Conn
	name string
}

var (
	newClients  chan client
	leftClients chan client
	messages    chan string
	namereg     *regexp.Regexp
)

func init() {
	newClients = make(chan client)
	leftClients = make(chan client)
	messages = make(chan string)
	namereg = regexp.MustCompile(`[a-z|A-Z|0-9]*`)
}

func main() {
	listener, err := net.Listen("tcp", ":8765")
	if err != nil {
		panic(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]struct{})

	for {
		select {
		case newCli := <-newClients:
			fmt.Println("Register new user:", newCli.name)
			clients[newCli] = struct{}{}
			pass := make([]string, 0, len(clients))
			for cli := range clients {
				pass = append(pass, cli.name)
			}
			res := "USERS" + strings.Join(pass, ",")
			for cli := range clients {
				cli.ch <- res
			}
			go func() {
				fmt.Println("Send to messages", newCli.name+" присоединился.")
				messages <- newCli.name + " присоединился."
			}()
		case leftCli := <-leftClients:
			delete(clients, leftCli)
			close(leftCli.ch)
			pass := make([]string, 0, len(clients))
			for cli := range clients {
				pass = append(pass, cli.name)
			}
			res := "USERS" + strings.Join(pass, ",")
			for cli := range clients {
				cli.ch <- res
			}
			go func() {
				messages <- leftCli.name + " вышел."
			}()
		case msg := <-messages:
			msg = "MESSAGE" + msg
			for cli := range clients {
				cli.ch <- msg
			}
		}
	}
}

func handleConn(conn net.Conn) {
	msgEnd := "\xe2\x90\x9c"

	fmt.Println("New connection aquired" +
		conn.RemoteAddr().String())
	cli := client{}
	cli.ch = make(chan string)
	cli.conn = conn
	defer conn.Close()

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading name")
		close(cli.ch)
		return
	}
	nameMsg := string(buff[:n])
	if len(nameMsg) < 5 ||
		nameMsg[:4] != "NAME" ||
		!namereg.Match([]byte(nameMsg[4:])) {
		fmt.Println("Got bad name format:" + nameMsg)
		close(cli.ch)
		return
	}
	cli.name = nameMsg[4:]
	newClients <- cli

	go func() {
		for msg := range cli.ch {
			fmt.Println("Sending message to", cli.name, msg)
			_, err := conn.Write([]byte(msg + msgEnd))
			if err != nil {
				fmt.Printf("%s (%s): %v\n",
					conn.RemoteAddr().String(), cli.name, err)
			}
		}
	}()

	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println(conn.RemoteAddr(), err)
			break
		}
		fmt.Println("Message got from", cli.name+":", string(buff[:n]))
		messages <- fmt.Sprintf("%s: %s",
			cli.name, string(buff[:n]))
	}
	leftClients <- cli
}
