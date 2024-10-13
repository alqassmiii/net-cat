package netserver

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	Msgs       chan clientMessage
	clients    map[net.Conn]client
	logFile    *os.File // File for logging messages
	oldMsgs    []string // Slice to hold old messages

}

// Struct to hold the message and the sender connection
type clientMessage struct {
	Message string
	Sender  net.Conn
}

// Create a new server with a listener address
func NewServer(listenAddr string) (*server, error) {
	// Open or create the log file
	os.Remove("chat_logs.txt")
	logFile, err := os.OpenFile("chat_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %v", err)
	}

	return &server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		Msgs:       make(chan clientMessage, 10),
		clients:    make(map[net.Conn]client),
		logFile:    logFile,
	}, nil
}

// Start the server and listen for incoming connections
func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	defer s.logFile.Close() // Ensure log file is closed when the server stops
	s.ln = ln

	go s.accept()

	<-s.quitch
	close(s.Msgs)

	return nil
}

// Log a message to both the console and the log file
func (s *server) LogMessage(message string) {
	// Log to the console
	fmt.Println(message)
	// Write to the log file
	if _, err := s.logFile.WriteString(message + "\n"); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
	// Ensure the message is flushed to disk
	s.logFile.Sync()
}

// Accept incoming client connections
func (s *server) accept() {
	for {
		if len(s.clients) >= 3 {
			fmt.Println("Connection limit reached, rejecting new connections.")
			continue
		}

		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		if len(s.clients) >= 3 { // Double-check after accepting
			conn.Write([]byte("Server is full, please try again later.\n"))
			conn.Close()
			continue
		}

		fmt.Println("Listening on the port", s.listenAddr)
		go s.handleConnection(conn)
	}
}

// Broadcast messages to all connected clients except the sender
func (s *server) Broadcast(message string, sender net.Conn) {
	s.oldMsgs = append(s.oldMsgs, message)
	for conn, client := range s.clients {
		if conn != sender {
			client.conn.Write([]byte(message + "\n"))
		}
	}
}

func (s *server) read(conn net.Conn, clientInfo client) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			exitMessage := fmt.Sprintf("%s user left: %s", clientInfo.name, err)
			s.LogMessage(exitMessage) // Log the disconnection
			s.Broadcast(fmt.Sprintf("%s has left our chat...", clientInfo.name), conn)
			delete(s.clients, conn)
			return
		}
		message = strings.TrimSpace(message)

		timestamp := time.Now().Format("2006-01-02 15:04:05")

		formattedMessage := fmt.Sprintf("[%s][%s]: %s", timestamp, clientInfo.name, message)

		s.Msgs <- clientMessage{
			Message: formattedMessage,
			Sender:  conn,
		}
	}
}

// Handle each client connection
func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to TCP-Chat!\n"))

	logo, err := os.ReadFile("linuxlogo.txt")
	if err != nil {
		fmt.Println("Error reading logo file:", err)
		return
	}
	conn.Write(logo)
	conn.Write([]byte("\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	name = strings.TrimSpace(name)

	for _, clientInfo := range s.clients {
		if clientInfo.name == name {
			conn.Write([]byte("Sorry! The name you are trying to enter is already in use.\nPlease try to use another name."))
			conn.Close()
			return
		}
	}

	clientInfo := client{
		name: name,
		from: conn.RemoteAddr().String(),
		conn: conn,
	}

	s.clients[conn] = clientInfo

	joinMessage := fmt.Sprintf("Client %s  connected", clientInfo.name)
	s.LogMessage(joinMessage) // Log the connection event

	s.Broadcast(fmt.Sprintf("%s has joined our chat...", clientInfo.name), conn)

	for _, msg := range s.oldMsgs {
		conn.Write([]byte(msg + "\n"))
	}

	s.read(conn, clientInfo)
}
