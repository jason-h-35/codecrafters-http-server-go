package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ListenAddress = "0.0.0.0:4221"
)

type HTTPRequest struct {
	method    string
	resource  string
	version   string
	host      string
	userAgent string
	acceptEnc string
}

func NewHTTPRequest(buf []byte) HTTPRequest {
	s := string(buf)
	lines := strings.Split(s, "\r\n")
	methodLine := strings.Split(lines[0], " ")
	var host, userAgent, acceptEnc string
	for _, l := range lines[1:] {
		if strings.HasPrefix(l, "Host:") {
			host = strings.Split(l, " ")[1]
		}
		if strings.HasPrefix(l, "User-Agent:") {
			userAgent = strings.Split(l, " ")[1]
		}
		if strings.HasPrefix(l, "Accept-Encoding:") {
			acceptEnc = strings.Split(l, " ")[1]
		}
	}
	return HTTPRequest{
		method:    methodLine[0],
		resource:  methodLine[1],
		version:   methodLine[2],
		host:      host,
		userAgent: userAgent,
		acceptEnc: acceptEnc,
	}
}

func handleRequest(requestBuf []byte) []byte {
	request := NewHTTPRequest(requestBuf)
	fmt.Printf("%s", requestBuf)
	var response []byte
	if strings.Compare(request.resource, "/") == 0 {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(request.resource, "/echo/") {
		body := []byte(request.resource[6:])
		headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n", len(body))
		header := []byte(headerString)
		response = append(header, body...)
	} else if strings.Compare(request.resource, "/user-agent") == 0 {
		body := []byte(request.userAgent)
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
