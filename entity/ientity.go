package gtentity

import "github.com/gtechx/base/gtnet"

type IEntity interface {
	ID() uint64
	UID() uint64
	APPID() uint64
	ZONE() uint32
	Conn() gtnet.IConn

	ForceOffline()
	RPC(firstmsgid uint8, secondmsgid uint8, params ...interface{})
}
