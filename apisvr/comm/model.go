package comm

type CommHandler interface {
	Connect() (bool, error)
	Send(b []byte) (int, error)
	Read(data []byte) (int, error)
	IsConnected() bool
	Close()
}

func NewTcpHandler(name string, port int, addr string) TcpHandler {
	return newTcpHandler(name, port, addr)
}

// func NewTcpHandler2(name string, port int, addr string) TcpHandler2 {
// 	return newTcpHandler2(name, port, addr)
// }
