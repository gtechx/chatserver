package main

import "net"

type Sess struct {
	id   uint64
	conn net.Conn
}

func (s *Sess) start() {

}
