package main

import "net"

type ISession interface {
	ID() uint64
	AppName() string
	ZoneName() string
	Send(interface{})
}

type Sess struct {
	appname  string
	zonename string
	id       uint64
	conn     net.Conn
}

func (s *Sess) ID() uint64 {
	return s.id
}

func (s *Sess) AppName() string {
	return s.appname
}

func (s *Sess) ZoneName() string {
	return s.zonename
}

func (s *Sess) start() {

}

func (s *Sess) Send(interface{}) {

}
