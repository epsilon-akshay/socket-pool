package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	// Connect to the server on port 8080.
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Get()
	// Write data to the connection.
	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read data from the connection.
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data that was read from the connection.
	fmt.Println(string(buf[:n]))
	conn.Close()
}
