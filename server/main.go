package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Accept an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go process(conn)
	}
}

func process(conn net.Conn) {

	response := make([]byte, 0, 4096) // Initialize an empty byte slice to hold the response

	buffer := make([]byte, 1044) // Create a buffer for reading
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Reached the end of the response
			}
			fmt.Println("Error reading response:", err)
			return
		}
		response = append(response, buffer[:n]...)
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
		fmt.Println("wrote respipnse", string(buffer))
	}

}
func process2(conn net.Conn) {
	response := make([]byte, 0, 4096)
	for {
		buffer, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("breakkkkkkk IOF")
				break // Reached the end of the response
			}

			fmt.Println("Error reading response:", err)
			return

		}
		response = append(response, buffer...)
		fmt.Println("read this much")
		_, err = conn.Write([]byte("success"))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}
