package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ListenAddress = "0.0.0.0:4221"
)

var DirFlag *string

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
		if p := "host:"; strings.HasPrefix(strings.ToLower(l), p) {
			host = strings.TrimSpace(l[len(p):])
		}
		if p := "user-agent:"; strings.HasPrefix(strings.ToLower(l), p) {
			userAgent = strings.TrimSpace(l[len(p):])
		}
		if p := "accept-encoding:"; strings.HasPrefix(strings.ToLower(l), p) {
			acceptEnc = strings.TrimSpace(l[len(p):])
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
	// fmt.Printf("%s", requestBuf)
	fmt.Printf("%s\n\n", request)
	var response []byte
	if request.method == "GET" {
		if request.resource == "/" {
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		} else if strings.HasPrefix(request.resource, "/echo/") {
			body := []byte(request.resource[6:])
			fmt.Println(body)
			headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n", len(body))
			header := []byte(headerString)
			response = append(header, body...)
		} else if request.resource == "/user-agent" {
			body := []byte(request.userAgent)
			headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n", len(body))
			header := []byte(headerString)
			response = append(header, body...)
		} else if *DirFlag != "" && strings.HasPrefix(request.resource, "/files/") {
			tokens := strings.Split(request.resource, "/")[2:] //everything after /files/
			subpath := strings.Join(tokens, "/")
			path := *DirFlag + subpath
			fileContents, err := os.ReadFile(path)
			if err != nil {
				return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
			}
			body := []byte(fileContents)
			headerString := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n", len(body))
			header := []byte(headerString)
			response = append(header, body...)
		} else {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		}
	} else if request.method == "POST" {
		if strings.HasPrefix(request.resource, "/files/") {
			tokens := strings.Split(request.resource, "/")[2:]
			subpath := strings.Join(tokens, "/")
			path := *DirFlag + subpath
			err := os.WriteFile(path, []byte(request.resource), 0644)
			if err != nil {
				fmt.Println("Error writing file: ", err.Error())
			}
			response = []byte("HTTP/1.1 201 OK\r\n\r\n")
		}
	}
	return append(response, []byte("\r\n")...)
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
	DirFlag = flag.String("directory", "", "serve a directory from local filesystem over HTTP")
	flag.Parse()
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
