package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn     net.Conn
	username string
}

const maxClients = 10

var (
	clients  = make(map[net.Conn]Client) // Stores connected clients
	messages []string // Stores chat history
	mutex    sync.Mutex  // Mutex to ensure safe concurrent access

)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// Greet client and prompt for name
	// Send welcome message and get client name
	welcomeMessage := "Welcome to TCP-Chat!\n" +
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
		" |    .       | ' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     -'       --'\n"
	conn.Write([]byte(welcomeMessage + "Welcome to TCP-Chat!\n[ENTER YOUR NAME]: "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Ensure a valid name
	if name == "" {
		conn.Write([]byte("Name cannot be empty.\n"))
		return
	}

	// Check if server is full
	mutex.Lock()
	if len(clients) >= maxClients {
		conn.Write([]byte("Server is full. Try again later.\n"))
		mutex.Unlock()
		return
	}

	// Add client to the list
	clients[conn] = Client{conn: conn, username: name}
	mutex.Unlock()

	// Notify others of the new connection
	join := fmt.Sprintf("%s has joined the chat", name)
	// Send chat history to the new client
	for _, msg := range messages {
		conn.Write([]byte(msg + "\n"))
	}

	BroadcastMessage(join, conn)

	// Listen for messages from this client
	for {
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), name)))
		message, err := reader.ReadString('\n')
		if err != nil {
			break // Client disconnected
		}

		message = strings.TrimSpace(message)
		if message != "" {
			BroadcastMessage(fmt.Sprintf("[%s][%s]:%s", time.Now().Format("2006-01-02 15:04:05"), name, message), conn)
		}
	}

	// Notify others when the client leaves
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	BroadcastMessage(fmt.Sprintf("%s has left the chat", name), conn)
}

// Broadcast a message to all clients except the sender
func BroadcastMessage(message string, exclude net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	messages = append(messages, message) // Store the message

	for conn, client := range clients {
		if conn != exclude { // Don't send the message to the sender
			prompt := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), client.username)
			conn.Write([]byte("\n" + message + "\n" + prompt))
			// Prompt the user again after a message
		}
	}
}