package main

import (
	"fmt"
	netserver "netcat/server"
	"os"
	"strconv"
)

func main() {
	var port string
	if len(os.Args) == 1 {
		port = ":8989" // Default port
	} else if len(os.Args) == 2 && len(os.Args[1]) <= 4 {
		if _, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println("Please enter a valid numeric port")
			fmt.Println("Usage: go run main.go [Numeric port]")
			return
		}
		port = ":" + os.Args[1] // Get the port from the command line
	} else {
		fmt.Println("Usage: go run main.go [Numeric port]") 
		return
	}
	fmt.Println("Server is running on Port", port)
	server, err := netserver.NewServer(port, 10) // Create a new server and set the maximum number of clients
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}

	go func() { // Start a goroutine to handle messages
		for msg := range server.Msgs {
			server.LogMessage(msg.Message)           
			server.Broadcast(msg.Message, msg.Sender) 
		}
	}()

	err = server.Start() // Start the server
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
