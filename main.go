package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type client struct {
	name    string
	from    string
	payload string
}

type server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgs       chan client
}

// Create a new server with a listener address
func NewServer(listenAddr string) *server {
	return &server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgs:       make(chan client, 10),
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

		fmt.Println("Listening on thr port", s.listenAddr)
		go s.handleConnection(conn) // Handle client in a goroutine
	}
}
func (s *server) read(conn net.Conn, clientInfo client) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n') // Read message until newline
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		message = strings.TrimSpace(message)

		// Get the current timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		// Send the message to the main server loop with client info and timestamp
		s.msgs <- client{
			name:    clientInfo.name,
			from:    clientInfo.from,
			payload: fmt.Sprintf("[%s][%s]: %s", timestamp, clientInfo.name, message),
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
	}

	fmt.Printf("Client %s (%s) connected\n", clientInfo.name, clientInfo.from )

	// Start reading messages from the client
	s.read(conn, clientInfo)
}

// Read messages from the client

func main() {
	server := NewServer(":8989")

	// Goroutine to print received messages from clients
	go func() {
		for msg := range server.msgs {
			// Print the message in the requested format
			fmt.Println(msg.payload)
		}
	}()

	// Start the server
	err := server.Start()
	if err != nil {
		fmt.Println("Server error:", err)
	}
}