package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// listen
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	// conn from listener
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	// send response over conn.
	response := []byte("HTTP/1.1 200 OK\r\n\r\n")
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
