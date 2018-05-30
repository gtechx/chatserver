package main

import (
	"fmt"
	"net"
)

type ISession interface {
	ID() uint64
	Account() string
	AppName() string
	ZoneName() string
	Send(interface{})
	Start()
	Stop()
}

type Sess struct {
	account  string
	appname  string
	zonename string
	id       uint64
	conn     net.Conn

	sendChan chan []byte
	quitChan chan int
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

func (s *Sess) Start() {
	s.quitChan = make(chan int, 1)
	s.sendChan = make(chan []byte, 2)
	go s.startRecv()
	startSend()
}

func (s *Sess) Stop() {
	s.quitChan <- 1
}

func (s *Sess) Send(buff []byte) {
	s.sendChan <- buff
}

func (s *Sess) startRecv() {
	for {
		msgtype, id, size, msgid, databuff, err := readMsgHeader(s.conn)
		if err != nil {
			fmt.Println("readMsgHeader error:" + err.Error())
			break
		}
		switch msgtype {
		case TickFrame:
		case EchoFrame:
			senddata := packageMsg(EchoFrame, id, msgid, databuff)
			s.sendChan <- senddata
		default:
			errcode, ret := HandleMsg(msgid, s, databuff)
			if errcode != ERR_MSG_INVALID {
				senddata := packageMsg(RetFrame, id, msgid, ret)
				s.sendChan <- senddata
			}
		}
	}
	s.quitChan <- 1
}

func (s *Sess) startSend() {
	for {
		select {
		case <-quitChan:
			count := len(sendChan)
			for i := 0; i < count; i++ {
				databuff := <-sendChan
				_, err := s.conn.Write(databuff)
				if err != nil {
					fmt.Println("err Send:" + err.Error())
					break
				}
			}
			break
		case databuff := <-sendChan:
			_, err := s.conn.Write(databuff)
			if err != nil {
				fmt.Println("err Send:" + err.Error())
				// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				// 	//srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				// 	//time.Sleep(tempDelay)
				// 	continue
				// }
				break
			}
		}
	}
	fmt.Println("remove session from sessmgr..")
	SessMgr().DelSess(s.id)
}
