package main

import (
	"fmt"
	netserver "netcat/server"
	"os"
)

func main() {
	var port string
	if len(os.Args) == 1 {
		port = ":8989"
	} else {
		port = ":" + os.Args[1]
	}
	fmt.Println("Server is running... Port:", port)
	server, err := netserver.NewServer(port,10)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	go func() {
		for msg := range server.Msgs {
			server.LogMessage(msg.Message) // Log the message
			server.Broadcast(msg.Message, msg.Sender)
		}
	}()

	err = server.Start()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
