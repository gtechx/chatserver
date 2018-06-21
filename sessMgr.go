package main

import (
	"net"
	"sync"
)

type SessManager struct {
	sessMap *sync.Map
}

var sessmgr *SessManager

func SessMgr() *SessManager {
	if sessmgr == nil {
		sessmgr = &SessManager{sessMap: &sync.Map{}}
	}
	return sessmgr
}

func (sm *SessManager) CreateSess(conn net.Conn, appname, zonename, account string, id uint64) ISession {
	sess := &Sess{account: account, appname: appname, zonename: zonename, id: id, conn: conn}
	sesslist := sm.GetSess(id)
	if sesslist == nil {
		sesslist = map[string]ISession{}
		sm.sessMap.Store(id, sesslist)
	}
	oldsess, ok := sesslist[appname]
	if ok {
		oldsess.KickOut()
	}
	sesslist[appname] = sess
	return sess
}

func (sm *SessManager) DelSess(id uint64) {
	sm.sessMap.Delete(id)
}

func (sm *SessManager) GetSess(id uint64) map[string]ISession {
	sesslist, ok := sm.sessMap.Load(id)
	if ok {
		return sesslist.(map[string]ISession)
	}
	return nil
}

func (sm *SessManager) SendMsgToId(id uint64, msg []byte) bool {
	sesslist := sm.GetSess(id)
	if sesslist != nil {
		for _, sess := range sesslist {
			sess.(*Sess).Send(msg)
		}
		return true
	}
	return false
}
