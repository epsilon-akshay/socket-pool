package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Pool struct {
	connections map[int]net.Conn
	idleConnMu  sync.Mutex
	idleConn    chan int
}

func pop(idleConn chan int) int {
	return <-idleConn
}

func push(idleConn chan int, id int) {
	idleConn <- id
}

func main() {
	// Connect to the server on port 8080.
	pool := &Pool{}
	pool.idleConnMu = sync.Mutex{}
	pool.connections = make(map[int]net.Conn)
	pool.idleConn = make(chan int, 3)
	conn, err := net.Dial("tcp", "localhost:8080")
	conn2, err := net.Dial("tcp", "localhost:8080")
	conn3, err := net.Dial("tcp", "localhost:8080")

	pool.connections[1] = conn
	pool.connections[2] = conn2
	pool.connections[3] = conn3

	setIdle(pool, 1)
	setIdle(pool, 2)
	setIdle(pool, 3)

	wg := &sync.WaitGroup{}

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		m := "hihi+" + strconv.Itoa(i)
		go connectAndWrite(m, wg, pool, i)
	}

	wg.Wait()
	conn.Close()
	conn2.Close()
	conn3.Close()
}

func get(pool *Pool) (net.Conn, int) {
	pool.idleConnMu.Lock()
	defer pool.idleConnMu.Unlock()
	lastConnId := pop(pool.idleConn)
	return pool.connections[lastConnId], lastConnId
}

func setIdle(pool *Pool, id int) {
	push(pool.idleConn, id)
}

func connectAndWrite(message string, wg *sync.WaitGroup, pool *Pool, id int) {
	conn, id := get(pool)
	//fmt.Println("connection used is", conn, "and id is", id)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 2048)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	setIdle(pool, id)

	//fmt.Println("for req", message, "response is", string(buf[:n]))
	wg.Done()

}

func connectAndWriteWithoutPool(message string, wg *sync.WaitGroup, conn net.Conn, id int) {
	//fmt.Println("connection used is", conn, "and id is", id)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 2048)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("for req", message, "response is", string(buf[:n]))
	wg.Done()
	conn.Close()
}
