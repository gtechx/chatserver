package gtservice

import (
	"time"

	"github.com/nature19862001/base/gtnet"
)

type Service struct {
	name      string
	net       string
	addr      string
	startTime int64

	connHandler func(gtnet.IConn)
	server      *gtnet.Server
}

func NewService(name string, net string, addr string, connhandler func(gtnet.IConn)) *Service {
	return &Service{name: name, net: net, addr: addr, connHandler: connhandler}
}

func (this *Service) Start() error {
	var err error

	server := gtnet.NewServer(this.net, this.addr, this.connHandler)

	err = server.Start()
	if err == nil {
		this.server = server
		this.startTime = time.Now().Unix()
	}

	return err
}

func (this *Service) Stop() error {
	return this.server.Stop()
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Net() string {
	return this.net
}

func (this *Service) Addr() string {
	return this.addr
}

func (this *Service) StartTime() int64 {
	return this.startTime
}
