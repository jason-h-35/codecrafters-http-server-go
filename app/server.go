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
		lines := strings.Split(request, "\r\n")
		startLine := strings.Split(lines[0], " ")
		// fmt.Println(startLine[1])
		var response []byte
		if strings.Compare(startLine[1], "/") == 0 {
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		} else {
			fmt.Println(startLine[1])
			body := []byte(startLine[1][1:])
			fmt.Println(body)
			header := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 3\r\n\r\n")
			response = append(header, body...)
		}
		// send response over conn.
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			os.Exit(1)
		}
	}
}
