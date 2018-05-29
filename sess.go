package main

import "net"

type ISession interface {
	ID() uint64
	Account() string
	AppName() string
	ZoneName() string
	Send(interface{})
}

type Sess struct {
	account  string
	appname  string
	zonename string
	id       uint64
	conn     net.Conn
}

func (s *Sess) ID() uint64 {
	return s.id
}

func (s *Sess) Account() string {
	return s.account
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
