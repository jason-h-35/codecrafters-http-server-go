package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const ListenAddress = "0.0.0.0:4221"

type HTTPRequest struct {
	method    string
	resource  string
	version   string
	host      string
	userAgent string
	accept    string
	bytes     []byte
}

func NewHTTPRequest(buf []byte) HTTPRequest {
	s := string(buf)
	lines := strings.Split(s, "\r\n")
	tokens := make([][]string, len(lines))
	for i := range lines {
		tokens[i] = strings.Split(lines[i], " ")
	}
	return HTTPRequest{
		method:    tokens[0][0],
		resource:  tokens[0][1],
		version:   tokens[0][2],
		host:      tokens[1][1],
		userAgent: tokens[2][1],
		accept:    tokens[3][1],
		bytes:     buf,
	}
}

func handleRequest(requestBuf []byte) []byte {
	request := string(requestBuf)
	fmt.Printf(request)
	lines := strings.Split(request, "\r\n")
	startLine := strings.Split(lines[0], " ")
	var response []byte

	if strings.Compare(startLine[1], "/") == 0 {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(startLine[1], "/echo/") {
		body := []byte(startLine[1][6:])
		headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n", len(body))
		header := []byte(headerString)
		response = append(header, body...)
	} else if strings.Compare(startLine[1], "/user-agent") == 0 {
		userAgent := strings.Split(lines[2], " ")[1]
		body := []byte(userAgent)
		headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n", len(body))
		header := []byte(headerString)
		response = append(header, body...)
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	return response
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
	}

	response := handleRequest(request[:n])
	// send response over conn.
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// listen
	l, err := net.Listen("tcp", ListenAddress)
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
		go handleConnection(conn)
	}
}
