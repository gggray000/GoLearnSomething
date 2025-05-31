package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	serverHost = flag.String("host", "localhost", "server host")
	serverPort = flag.Int("port", 8972, "server port")
)

func main() {
	flag.Parse()

	// Connecting to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *serverHost, *serverPort))
	if err != nil {
		fmt.Printf("Error connecting to server: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Failed to connect to server")

	// go routine for accepting messages from server
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("\n Disconneted from server: %s\n", err)
				os.Exit(1)
			}
			fmt.Print(string(buffer[:n]))
		}

	}()

	// read message from stdin and send it
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		msg = strings.TrimSpace(msg)

		if msg == "" {
			continue
		}

		// sending messages to server
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Failed to write to server: ", err)
			break
		}

		// quit by entering "/quit"
		if msg == "/quit" {
			fmt.Println("Goodbye!")
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading stdin:", err)
	}

}
