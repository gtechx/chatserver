package main

import (
	"fmt"
	"net"

	"github.com/gtechx/base/collections"
	"github.com/gtechx/chatserver/db"
)

type ISession interface {
	ID() uint64
	Account() string
	AppName() string
	ZoneName() string
	NickName() string
	Platform() string
	Send(buff []byte) bool
	Start()
	Stop()
	KickOut()
}

type Sess struct {
	appdata  *gtdb.AppData
	conn     net.Conn
	platform string

	sendList  *collections.SafeList
	quitChan  chan bool
	errorChan chan bool
	isClosed  bool
}

func (s *Sess) ID() uint64 {
	return s.appdata.ID
}

func (s *Sess) Account() string {
	return s.appdata.Account
}

func (s *Sess) AppName() string {
	return s.appdata.Appname
}

func (s *Sess) ZoneName() string {
	return s.appdata.Zonename
}

func (s *Sess) NickName() string {
	return s.appdata.Nickname
}

func (s *Sess) Platform() string {
	return s.platform
}

func (s *Sess) Start() {
	s.quitChan = make(chan bool, 2) //多个goroutine有可能会quit，所以需要两个，防止阻塞某个goroutine
	s.errorChan = make(chan bool, 1)
	s.sendList = collections.NewSafeList()
	go s.startRecv()
	s.startSend()
}

func (s *Sess) Stop() {
	s.isClosed = true
	s.quitChan <- true
}

func (s *Sess) KickOut() {
	senddata := packageMsg(RetFrame, 0, MsgId_KickOut, nil)
	s.Send(senddata)
	s.Stop()
}

func (s *Sess) Send(buff []byte) bool {
	if s.isClosed {
		return false
	}
	s.sendList.Put(buff)
	// select {
	// case s.sendChan <- buff:
	// case <-time.After(time.Millisecond * 100):
	// 	return false
	// }
	return true
}

var tickdata = []byte{TickFrame}

func (s *Sess) startRecv() {
	for {
		msgtype, id, _, msgid, databuff, err := readMsgHeader(s.conn)
		if err != nil {
			fmt.Println("readMsgHeader error:" + err.Error())
			s.errorChan <- true
			break
		}
		switch msgtype {
		case TickFrame:
			//s.sendChan <- tickdata
			s.sendList.Put(tickdata)
		case EchoFrame:
			senddata := packageMsg(EchoFrame, id, msgid, databuff)
			//s.sendChan <- senddata
			s.sendList.Put(senddata)
		default:
			if msgid != MsgId_ReqQuitChat {
				errcode, ret := HandleMsg(msgid, s, databuff)
				if errcode == ERR_MSG_INVALID {
					fmt.Println("ERR_MSG_INVALID")
					s.errorChan <- true
					goto end
				}
				if ret != nil {
					senddata := packageMsg(RetFrame, id, msgid, ret)
					//s.sendChan <- senddata
					s.sendList.Put(senddata)
				}
			} else {
				s.Stop()
				goto end
			}
		}
	}
end:
	fmt.Println("sess recv end")
}

func (s *Sess) startSend() {
	for {
		select {
		case <-s.quitChan:
			fmt.Println("sess start quit...")
			// count := len(s.sendChan)
			// for i := 0; i < count; i++ {
			// 	databuff := <-s.sendChan
			// 	_, err := s.conn.Write(databuff)
			// 	if err != nil {
			// 		fmt.Println("err Send:" + err.Error())
			// 		goto end
			// 	}
			// }
			for {
				data, err := s.sendList.Pop()
				if err != nil {
					break
				}

				databuff := data.([]byte)
				_, err = s.conn.Write(databuff)
				if err != nil {
					fmt.Println("err Send:" + err.Error())
					goto end
				}
			}
			goto end
		case <-s.errorChan:
			goto end
		case <-s.sendList.C:
			for {
				data, err := s.sendList.Pop()
				if err != nil {
					break
				}

				databuff := data.([]byte)

				_, err = s.conn.Write(databuff)
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
	}
end:
	fmt.Println("remove session from sessmgr..")
	SessMgr().DelSess(s)
	//count := len(s.sendChan)
	for {
		data, err := s.sendList.Pop()
		if err != nil {
			break
		}

		databuff := data.([]byte)
		SessMgr().TrySaveOfflineMsg(s.ID(), databuff)
	}

	// for i := 0; i < count; i++ {
	// 	databuff := <-s.sendChan
	// 	SessMgr().TrySaveOfflineMsg(s.ID(), databuff)
	// }
}
