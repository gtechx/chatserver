package main

import (
	"net"
	"sync"
)

type SessMgr struct {
	sessMap *sync.Map
}

func (sm *SessMgr) CreateSess(conn net.Conn, uint64 id) *Sess {
	sess := &Sess{id, conn}
	sm.sessMap.Store(id, sess)
	return sess
}

func (sm *SessMgr) DelSess(uint64 id) {
	sm.sessMap.Delete(id)
}

func (sm *SessMgr) GetSess(uint64 id) *Sess {
	sess, ok := sm.sessMap.Load(id)
	if ok {
		return sess.(*Sess)
	}
	return nil
}
