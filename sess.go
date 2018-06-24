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
	NickName() string
	Send(buff []byte)
	Start()
	Stop()
	KickOut()
}

type Sess struct {
	account  string
	appname  string
	zonename string
	nickname string
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

func (s *Sess) NickName() string {
	return s.nickname
}

func (s *Sess) Start() {
	s.quitChan = make(chan int, 1)
	s.sendChan = make(chan []byte, 2)
	go s.startRecv()
	s.startSend()
}

func (s *Sess) Stop() {
	s.quitChan <- 1
}

func (s *Sess) KickOut() {
	senddata := packageMsg(RetFrame, 0, MsgId_KickOut, nil)
	s.Send(senddata)
	s.Stop()
}

func (s *Sess) Send(buff []byte) {
	s.sendChan <- buff
}

func (s *Sess) startRecv() {
	for {
		msgtype, id, _, msgid, databuff, err := readMsgHeader(s.conn)
		if err != nil {
			fmt.Println("readMsgHeader error:" + err.Error())
			break
		}
		switch msgtype {
		case TickFrame:
			s.sendChan <- []byte{TickFrame}
		case EchoFrame:
			senddata := packageMsg(EchoFrame, id, msgid, databuff)
			s.sendChan <- senddata
		default:
			if msgid != MsgId_ReqQuitChat {
				errcode, ret := HandleMsg(msgid, s, databuff)
				if errcode == ERR_MSG_INVALID {
					goto end
				}
				if ret != nil {
					senddata := packageMsg(RetFrame, id, msgid, ret)
					s.sendChan <- senddata
				}
			} else {
				goto end
			}
		}
	}
end:
	s.quitChan <- 1
	fmt.Println("sess recv end")
}

func (s *Sess) startSend() {
	for {
		select {
		case <-s.quitChan:
			fmt.Println("sess start quit...")
			count := len(s.sendChan)
			for i := 0; i < count; i++ {
				databuff := <-s.sendChan
				_, err := s.conn.Write(databuff)
				if err != nil {
					fmt.Println("err Send:" + err.Error())
					goto end
				}
			}
			goto end
		case databuff := <-s.sendChan:
			_, err := s.conn.Write(databuff)
			if err != nil {
				fmt.Println("err Send:" + err.Error())
				// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
				// 	//srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				// 	//time.Sleep(tempDelay)
				// 	continue
				// }
				goto end
			}
		}
	}
end:
	fmt.Println("remove session from sessmgr..")
	SessMgr().DelSess(s.id)
}
