package netserver

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
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
	logFile    *os.File   // File for logging messages
	oldMsgs    []string   // Slice to hold old messages
	maxClients int        // Maximum number of clients allowed
	mu         sync.Mutex // Mutex for locking and unlocking
}

// Struct to hold the message and the sender connection
type clientMessage struct {
	Message string
	Sender  net.Conn
}

// Create a new server with a listener address
func NewServer(listenAddr string, maxClients int) (*server, error) {
	// Open or create the log file
	os.Remove("chat_logs.txt")
	logFile, err := os.OpenFile("chat_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %v", err)
	}

	s := &server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		Msgs:       make(chan clientMessage),
		clients:    make(map[net.Conn]client),
		logFile:    logFile,
		maxClients: maxClients,
	}

	return s, nil
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
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go s.handleNewConnection(conn)
	}
}

// Function to handle new connections
func (s *server) handleNewConnection(conn net.Conn) {
	s.mu.Lock()
	if len(s.clients) >= s.maxClients {
		conn.Write([]byte("Maximum connection limit reached. Please try again later...\n"))
		conn.Close()
		s.mu.Unlock()
		fmt.Println("Max clients reached. Connection refused.")
		return
	}
	s.clients[conn] = client{conn: conn}
	s.mu.Unlock()

	// Handle the client connection
	go s.handleClient(conn)
}

// Function to handle client communication
func (s *server) handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to TCP-Chat!\n"))

	logo, err := os.ReadFile("linuxlogo.txt")
	if err != nil {
		fmt.Println("Error reading logo file:", err)
		return
	}
	conn.Write(logo)
	conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	name = strings.TrimSpace(name)

	s.mu.Lock()
	for _, clientInfo := range s.clients {
		if name == "" {
			conn.Write([]byte("Sorry! Empty name cant be accepted.\nDisconnected...\n"))
			conn.Close()
			s.mu.Unlock()
			return
		} else if clientInfo.name == name {
			conn.Write([]byte("Sorry! The name you are trying to enter is already in use.\nDisconnected...\n"))
			conn.Close()
			s.mu.Unlock()
			return
		}

	}

	clientInfo := client{
		name: name,
		from: conn.RemoteAddr().String(),
		conn: conn,
	}

	s.clients[conn] = clientInfo
	s.mu.Unlock()

	joinMessage := fmt.Sprintf("Client %s  connected", clientInfo.name)
	s.LogMessage(joinMessage) // Log the connection event

	s.Broadcast(fmt.Sprintf("%s has joined our chat...", clientInfo.name), conn)

	for _, msg := range s.oldMsgs {
		conn.Write([]byte(msg + "\n"))
	}

	s.read(conn, clientInfo)
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
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
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
