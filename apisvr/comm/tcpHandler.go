package comm

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Implement CommHandler
// ---------------------------------------------------------------------------
type TcpHandler struct {
	name      string
	port      int
	addr      string
	Connected bool
	//tcp       *net.TCPConn
	tcp net.Conn
}

var tcpMutex = &sync.Mutex{}

func newTcpHandler(name string, port int, addr string) TcpHandler {

	return TcpHandler{
		name: name,
		port: port,
		addr: addr,

		Connected: false,
	}
}

// ---------------------------------------------------------------------------
// SetCommEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) SetCommEnv(name string, port int, addr string) {
	t.name = name
	t.port = port
	t.addr = addr
}

// ---------------------------------------------------------------------------
// GetCommEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) GetCommEnv() (string, int) {

	return t.addr, t.port
}

// ---------------------------------------------------------------------------
// GetAddress
// ---------------------------------------------------------------------------
func (t *TcpHandler) GetAddress() string {
	return t.addr
}

// ---------------------------------------------------------------------------
// GetPort
// ---------------------------------------------------------------------------
func (t *TcpHandler) GetPort() int {
	return t.port
}

// ---------------------------------------------------------------------------
// Connect
// ---------------------------------------------------------------------------
func (t *TcpHandler) Connect() (bool, error) {
	target := fmt.Sprintf("%s:%d", t.addr, t.port)

	//log.Printf("Connect ... : %v", target)

	// raddr, err := net.ResolveTCPAddr("tcp", target)
	// if err != nil {
	// 	log.Printf("Connect err(1) : %v", err)
	// 	log.Fatal(err)
	// 	return false, err
	// }

	// connect to server
	//log.Printf("Connect ...(1) : %v", target)
	//conn, err := net.DialTCP("tcp", nil, raddr)
	conn, err := net.DialTimeout("tcp", target, time.Second*1)
	//log.Printf("Connect ...(2) : %v", target)
	if err != nil {
		//log.Fatal(err)
		//log.Printf("Connect err(2) : %v", err)
		return false, err
	}

	log.Printf("Connect ok : [%s]", target)

	t.Connected = true
	t.tcp = conn
	//t.tcp.SetNoDelay(false)

	return true, err
}

// ---------------------------------------------------------------------------
// SendMessage
// ---------------------------------------------------------------------------
func (t *TcpHandler) Send(b []byte) (int, error) {
	var cnt int = 0
	cnt, err := t.tcp.Write(b)
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}

// ---------------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------------
func (t *TcpHandler) Read(data []byte) (int, error) {
	var cnt int = 0
	cnt, err := t.tcp.Read(data)

	if err != nil {
		return cnt, err
	}

	return cnt, nil
}

func (t *TcpHandler) IsConnected() bool {
	var connected bool = false
	tcpMutex.Lock()
	connected = t.Connected
	tcpMutex.Unlock()

	return connected
}

// ---------------------------------------------------------------------------
// ClearEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) Close() {
	// clear tcp env
	// close tcp connect
	if !t.Connected {
		fmt.Printf("Aleady Close ok (%v)\n", t.Connected)
		return
	}

	t.Connected = false
	err := t.tcp.Close()
	if err != nil {
		fmt.Println("Close err..", err)
	}

	fmt.Printf("Close ok (%v)\n", t.Connected)
}
