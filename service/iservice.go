package gtservice

type IService interface {
	Start() error
	Stop() error
	Name() string
	Net() string
	Addr() string
	StartTime() int64
}
