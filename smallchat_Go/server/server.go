package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	maxClients        = 10
	maxNicknameLength = 32
)

var (
	serverPort = flag.Int("port", 8972, "server port")
)

type Client struct {
	conn     net.Conn
	nickname string
	readChan chan string
}

func (c Client) startRecovery() {
	for msg := range c.readChan {
		_, err := c.conn.Write([]byte(msg + "\r\n"))
		if err != nil {
			fmt.Println("Error sending recovery", err)
		}
	}
}

type ChatState struct {
	listener    net.Listener
	clientsLock sync.Mutex
	clients     map[net.Conn]*Client
	numClients  int
}

var chatState = &ChatState{
	clients: make(map[net.Conn]*Client),
}

func initChat() {
	var err error
	chatState.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		fmt.Println("Error listening on port", *serverPort, err)
		os.Exit(1)
	}
}

func broadcastMsg(client *Client, msg string) {
	chatState.clientsLock.Lock()
	for _, cli := range chatState.clients {
		if cli != client {
			cli.readChan <- ">> " + client.nickname + ": " + msg + "\r\n"
		}
	}
	chatState.clientsLock.Unlock()
}

func closeClient(client *Client) {
	chatState.clientsLock.Lock()
	close(client.readChan)
	err := client.conn.Close()
	if err != nil {
		fmt.Println("Error closing client", err)
		return
	}
	delete(chatState.clients, client.conn)
	chatState.numClients--
	chatState.clientsLock.Unlock()
}

func handleNewClient(client *Client) {
	welcomeMsg := "Welcome to Simple Chat!\n" +
		"Use /nickname to change nick name.\n" +
		"Use /cow to say something like a cow.\n" +
		"Use /dragon to say something like a dragon.\n"
	client.conn.Write([]byte(welcomeMsg))

	buffer := make([]byte, 256)
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			fmt.Printf("Client left: %s\n", client.conn.RemoteAddr().String())
			closeClient(client)
			return
		}

		msg := string(buffer[:n])
		msg = strings.TrimSpace(msg)
		if len(msg) > 0 && msg[0] == '/' {
			// Handling command lines
			parts := strings.SplitN(msg, " ", 2)
			cmd := parts[0]
			switch cmd {

			case "/nickname":
				if len(parts) > 1 {
					if len(parts[1]) > maxNicknameLength {
						client.conn.Write([]byte("Nickname is too long, please try again.\n"))
						continue
					}
					client.nickname = parts[1]
				}
				continue

			case "/cow":
				if len(parts) > 1 {
					out, err := exec.Command("cowsay", parts[1]).Output()
					if err != nil {
						fmt.Println("Error calling cowsay:", err)
						continue
					}
					broadcastMsg(client, "\n" + string(out))
				}
				continue

			case "/dragon":
				if len(parts) > 1 {
					out, err := exec.Command("cowsay", "-f", "dragon", parts[1]).Output()
					if err != nil {
						fmt.Println("Error calling cowsay:", err)
						continue
					}
					broadcastMsg(client, "\n" + string(out))
				}
				continue

			}
			// if cmd == "/nickname" && len(parts) > 1 {
			// 	if len(parts[1]) > maxNicknameLength {
			// 		client.conn.Write([]byte("Nickname is too long, please try again.\n"))
			// 		continue
			// 	}
			// 	client.nickname = parts[1]
			// }
			// continue
		}

		if len(msg) == 0 {
			continue
		}

		if buffer[0] == 253 || strings.ToLower(msg) == "quit" {
			closeClient(client)
			return
		}

		fmt.Printf("%s:%s\n", client.nickname, msg)

		broadcastMsg(client, msg)
	}
}

func main() {
	flag.Parse()
	initChat()
	for {
		conn, err := chatState.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}

		client := &Client{conn: conn}
		client.readChan = make(chan string, 5)

		chatState.clientsLock.Lock()
		if chatState.numClients >= maxClients {
			fmt.Printf("Max clients reached: %d, rejecting %s.\n", maxClients, client.conn.RemoteAddr().String())
			conn.Close()
			chatState.clientsLock.Unlock()
			continue
		}
		chatState.clients[conn] = client
		chatState.numClients++
		chatState.clientsLock.Unlock()

		go handleNewClient(client)
		go client.startRecovery()
		fmt.Printf("New client: %s\n", client.conn.RemoteAddr().String())
	}

}
