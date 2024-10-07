package main

import (
	"fmt"
	"net"
)

type client struct {
	from    string
    payload []byte
} 

type server struct {
	listenAddr string
	ln     net.Listener

	quitch chan struct{}

	msgs  chan client
}


func NewServer(listenAddr string) *server {
	return &server{
		listenAddr: listenAddr,
		quitch:make(chan struct{}),
		msgs: make(chan client, 10),
	}
}
func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.accept()

	<-s.quitch
	close(s.msgs)

	
	return nil

	
}

func (s *server) accept()  {
	for{
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("accept conn:", conn.RemoteAddr())
		go s.read(conn)
	}
}
func (s *server) read(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for{
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		s.msgs <- client{
			from : conn.RemoteAddr().String(),
			payload: buf[:n],
		} 
	}
	
}
func main (){
	server := NewServer(":8989")

	go func() {
		for msg := range server.msgs {
			fmt.Println("received  message from ",msg.from,":", string(msg.payload))
		}
	}()
	server.Start()
}