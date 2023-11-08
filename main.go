package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type Pool struct {
	connections map[int]*WrapperConn
	mutex       sync.Mutex
	idleConn    chan int
	maxConn     int
	connTimeout time.Duration
}

type WrapperConn struct {
	id         int
	Conn       net.Conn
	ConnReader *bufio.Reader
	ConnWriter *bufio.Writer
	IsActive   bool
}

func push(idleConn chan int, id int) {
	idleConn <- id
}

func main() {
	// Connect to the server on port 8080.
	pool := &Pool{}
	pool.mutex = sync.Mutex{}
	pool.maxConn = 3
	pool.connections = make(map[int]*WrapperConn)
	pool.idleConn = make(chan int, pool.maxConn)
	pool.connTimeout = 10 * time.Second

	wg := &sync.WaitGroup{}
	errChan := make(chan error)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		m := "hello " + strconv.Itoa(i)
		go connectAndWriteWithBuffers(m, wg, pool, errChan)
	}
	wg.Wait()
	for _, v := range pool.connections {
		v.Conn.Close()
	}
}

func waitOrGetIdle(pool *Pool) (*WrapperConn, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	select {
	case lastConnId := <-pool.idleConn:
		fmt.Println("found conn")
		pool.connections[lastConnId].IsActive = true
		return pool.connections[lastConnId], nil
	case <-time.After(pool.connTimeout):
		fmt.Println("oh nooo")
		return nil, errors.New("conn timeout")
	}
}

func addAndGet(pool *Pool) (*WrapperConn, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.connections) >= pool.maxConn {
		return nil, errors.New("max conn err")
	}

	newconn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	wc := &WrapperConn{
		id:         rand.Int(),
		Conn:       newconn,
		ConnReader: bufio.NewReader(newconn),
		ConnWriter: bufio.NewWriter(newconn),
		IsActive:   true,
	}
	pool.connections[wc.id] = wc
	return wc, nil
}

func setIdle(pool *Pool, id int) {
	//TODO: think of critical section here
	pool.connections[id].IsActive = false
	push(pool.idleConn, id)
}

func connectAndWriteWithoutPool(message string, wg *sync.WaitGroup, conn net.Conn, id int) {
	//fmt.Println("connection used is", conn, "and id is", id)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 10000)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(" id is", id, "response is", string(buf), " \n", "----------------------------")
	wg.Done()
	conn.Close()
}

func connectAndWriteWithBuffers(message string, wg *sync.WaitGroup, pool *Pool, errChan chan error) {

	wc, err := addAndGet(pool)
	if err != nil {
		wc, err = waitOrGetIdle(pool)
		if err != nil {
			errChan <- err
		}
	}

	_, err = wc.ConnWriter.Write([]byte(message))
	if err != nil {
		fmt.Println(err)
		errChan <- err
	}

	err = wc.ConnWriter.Flush()
	if err != nil {
		fmt.Println(err)
		errChan <- err
	}

	buf := make([]byte, 2048)
	_, err = wc.ConnReader.Read(buf)
	if err != nil {
		fmt.Println(err)
		errChan <- err
	}
	setIdle(pool, wc.id)

	fmt.Println(" id is", wc.id, "response is", string(buf), " \n", "----------------------------")
	wg.Done()
}
