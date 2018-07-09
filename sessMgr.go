package main

import (
	"net"
	"sync"

	"github.com/gtechx/chatserver/db"
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

func (sm *SessManager) CreateSess(conn net.Conn, tbl_appdata *gtdb.AppData, platform string) ISession {
	sess := &Sess{appdata: tbl_appdata, conn: conn, platform: platform}
	sesslist := sm.GetSess(tbl_appdata.ID)
	if sesslist == nil {
		sesslist = map[string]ISession{}
		sm.sessMap.Store(tbl_appdata.ID, sesslist)
	}
	oldsess, ok := sesslist[platform]
	if ok {
		oldsess.KickOut()
	}
	sesslist[platform] = sess
	return sess
}

func (sm *SessManager) DelSess(sess ISession) {
	sesslist := sm.GetSess(sess.ID())
	if sesslist != nil {
		delete(sesslist, sess.Platform())
	}
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
		flag := false
		for _, sess := range sesslist {
			flag = flag || sess.(*Sess).Send(msg)
		}
		return flag
	}
	return false
}

func (sm *SessManager) TrySaveOfflineMsg(id uint64, msg []byte) {
	sesslist := sm.GetSess(id)
	if len(sesslist) == 0 {
		gtdb.Manager().SendMsgToUserOffline(id, msg)
	}
}
