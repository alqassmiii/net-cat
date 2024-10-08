package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type client struct {
	name string
	from string
	conn net.Conn
}

type server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgs       chan string
	clients    map[net.Conn]client // Map to store all active clients
}

// Create a new server with a listener address
func NewServer(listenAddr string) *server {
	return &server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgs:       make(chan string, 10),
		clients:    make(map[net.Conn]client), // Initialize the clients map
	}
}

// Start the server and listen for incoming connections
func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.accept() // accept incoming connections

	<-s.quitch
	close(s.msgs)

	return nil
}

// Accept incoming client connections
func (s *server) accept() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("Listening on the port", s.listenAddr)
		go s.handleConnection(conn) // Handle client in a goroutine
	}
}

// Broadcast messages to all connected clients
func (s *server) broadcast(message string) {
	for _, client := range s.clients {
		client.conn.Write([]byte(message + "\n")) // Send message to each client
	}
}

func (s *server) read(conn net.Conn, clientInfo client) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n') // Read message until newline
		if err != nil {
			fmt.Printf("%s user left: %s", clientInfo.name, err)
			s.broadcast(fmt.Sprintf("%s has left our chat...", clientInfo.name))
			delete(s.clients, conn) // Remove the client from the active clients list
			return
		}
		message = strings.TrimSpace(message)

		// Get the current timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		if clientInfo.name == "Server"{
		// Create the message payload in the requested format
		formattedMessage := fmt.Sprintf("%s", message)
				// Send the message to all connected clients
				s.msgs <- formattedMessage
		}else {
		// Create the message payload in the requested format
		formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, clientInfo.name, message)
				// Send the message to all connected clients
				s.msgs <- formattedMessage
		}



	}
}

// Handle each client connection
func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Ask the client for their name
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte(
		    "         _nnnn_\n" +
			"        dGGGGMMb\n" +
			"       @p~qp~~qMb\n" +
			"       M|@||@) M|\n" +
			"       @,----.JM|\n" +
			"      JS^\\__/  qKL\n" +
			"     dZP        qKRb\n" +
			"    dZP          qKKb\n" +
			"   fZP            SMMb\n" +
			"   HZM            MMMM\n" +
			"   FqM            MMMM\n" +
			" __| \".        |\\dS\"qML\n" +
			" |    `.       | `' \\Zq\n" +
			"_)      \\.___.,|     .'\n" +
			"\\____   )MMMMMP|   .'\n" +
			"     `-'       `--'\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("error reading name:", err)
		return
	}
	name = strings.TrimSpace(name) // Remove trailing newline or spaces

	clientInfo := client{
		name: name,
		from: conn.RemoteAddr().String(),
		conn: conn, // Store the client's connection
	}

	// Add the new client to the active clients list
	s.clients[conn] = clientInfo

	fmt.Printf("Client %s (%s) connected\n", clientInfo.name, clientInfo.from)

	// Announce the new client to all other clients
	s.broadcast(fmt.Sprintf("%s has joined our chat", clientInfo.name))

	// Start reading messages from the client
	s.read(conn, clientInfo)
}

func main() {
	server := NewServer(":8989")

	// Goroutine to print and broadcast received messages from clients
	go func() {
		for msg := range server.msgs {
			// Print the message in the server console
			fmt.Println(msg)
			// Broadcast the message to all connected clients
			server.broadcast(msg)
		}
	}()

	// Start the server
	err := server.Start()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}