package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"handlers/functions"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on port:", port)
	for {
		connex, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handlers.HandleClient(connex)
	}
}
