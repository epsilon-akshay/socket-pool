package pooling

import "net"

type connection struct {
	conn   net.Conn
	active bool
}

func get() {

}
