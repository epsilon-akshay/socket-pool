package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe()
	// Accept an incoming connection.
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		newBuf := append([]byte("success"), buf...)
		fmt.Println(string(newBuf))
		_, err = conn.Write(newBuf)
		if err != nil {
			panic(err)
		}

	}
}
