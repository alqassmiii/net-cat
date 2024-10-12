# net-cat

Net-Cat is a simple TCP-based chat server implemented in Go. It allows multiple clients to join a chatroom, exchange messages, and keep logs of all chat activities. This project simulates the behavior of `NetCat` but extends it to include features like client names, message broadcasting, and chat history logging.

## Features

- **TCP Server**: Handles multiple client connections via TCP.
- **Client Naming**: Each client must provide a unique name upon connection.
- **Message Broadcasting**: Messages are broadcasted to all connected clients, excluding the sender.
- **Timestamped Messages**: All messages include a timestamp and client name.
- **Chat Logging**: All messages and events (join/leave) are logged to a file (`chat_logs.txt`).
- **Concurrency**: Supports multiple clients concurrently using goroutines.
- **Client Disconnection Handling**: Notifies all clients when a client disconnects.
- **Previous Message Upload**: New clients receive a history of past messages from the chat log.

## Project Structure
```t
net-cat/
├── client/
    └── client.go    # Client code
├── server/
    └── server.go    # Server code
├── main.go          # Main code
├── go.mod           # Go module file
├── LICENSE.md       # MIT License 
├── README.md        # README file
├── .gitignore       # Ignore OS binaries 
├── test.sh          # Test file
└── chat_logs.txt    # Chat log file
```
## Getting Started

### Prerequisites

- Go 1.20 or higher
- Terminal access

### Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/net-cat.git
```

2.	Navigate to the project directory:
```bash
cd net-cat
```


3.	Install dependencies:
```bash
go mod tidy
```


### Running the Server

1.	Start the server by specifying the address and port:
```bash
go run main.go
```
By default, the server will start on localhost:8989.

### Connecting Clients

Clients can connect to the server using a simple TCP client (like NetCat or a custom client).

Example using NetCat:
```bash
nc localhost <port>
```
Once connected, the client must provide a unique name before they can start chatting.

### Features in Action

	•	Broadcasting Messages: Once connected, all messages sent by a client will be broadcasted to all other clients.
	•	Logging: The server logs all chat activity, including client joins, leaves, and chat messages, to chat_logs.txt.
	•	Receiving Previous Messages: When a new client connects, they receive the last chat messages from the log file.
	•	Disconnection Handling: When a client disconnects, the server notifies all other clients.

## Example Usage

1.	Start the server for default port (8989):
```bash
go run main.go
```
or to add custom port:
```bash
go run main.go <port>
```

2.	Connect multiple clients:
```bash
nc localhost <port>
```

3.	Send messages and see them broadcasted across clients:
```bash
[2024-10-12 15:12:45][Alice]: Hello, everyone!
[2024-10-12 15:13:02][Bob]: Hi, Alice!
```

4.	When a client disconnects, the server notifies others:
```bash
Bob has left our chat...
```

## Limitations

	•	Client Limit: The server currently supports up to 10 clients. This can be changed in the code.
	•	No Encryption: Messages are sent in plain text over TCP. Consider adding SSL for secure communication in production environments.

## Future Improvements

	•	Encryption: Adding SSL/TLS to secure client-server communication.
	•	UI Enhancement: Implementing a terminal-based user interface (UI) for clients using gocui or a similar package.
	•	Client-Side Application: Create a Go-based or web-based client for a better user experience.

## Contributing

	1.	Fork the repository.
	2.	Create your feature branch (git checkout -b feature/NewFeature).
	3.	Commit your changes (git commit -m 'Add NewFeature').
	4.	Push to the branch (git push origin feature/NewFeature).
	5.	Open a Pull Request.

## Authors

	•	Sayed Ahmed Husain
	•	Qasim Aljaffer

## License

This project is licensed under the MIT License. See the LICENSE.md file for details.