package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkPool(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		m := "hihi+" + strconv.Itoa(i)
		go connectAndWrite(m, wg, pool, i)
	}
	wg.Wait()
	conn.Close()
	conn2.Close()
	conn3.Close()
}

func BenchmarkConn(b *testing.B) {

	wg := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			fmt.Println(err)
			return
		}

		m := "hihi+" + strconv.Itoa(i)
		go connectAndWriteWithoutPool(m, wg, conn, i)

	}
	wg.Wait()

}
