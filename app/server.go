package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// listen
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		// conn from listener
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		requestBuf := make([]byte, 1024)
		n, err := conn.Read(requestBuf)
		if err != nil {
			fmt.Println("Error reading connection: ", err.Error())
		}
		request := string(requestBuf[:n])
		fmt.Printf(request)
		lines := strings.Split(request, "\n")
		startLine := strings.Split(lines[0], " ")
		fmt.Println(startLine[1])
		var response []byte
		if strings.Compare(startLine[1], "/") == 0 {
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		} else {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		}
		// send response over conn.
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	}
}
