package netapp

import (
	"apisvr/comm"
	"time"

	dt "github.com/gdygd/goglib/datastructure"
)

// ---------------------------------------------------------------------------
// Protocol
// ---------------------------------------------------------------------------
const (
	PT_STX = 0x02
	PT_ETX = 0x03

	HEAD_SIZE = 5 // STX ~ OPCODE
	//TAIL_SIZE = 5 // CRC + ETX
	TAIL_SIZE = 1 // ETX

	PT_STX_IDX  = 0
	PT_LEN_IDX  = PT_STX_IDX + 1
	PT_SEQ_IDX  = PT_LEN_IDX + 2
	PT_OP_IDX   = PT_SEQ_IDX + 1
	PT_DATA_IDX = PT_OP_IDX + 1
)

const (
	PT_ACK  = 0x90 // test protocol
	PT_TEST = 0x99 // test protocol
)

// ---------------------------------------------------------------------------
// Rx State
// ---------------------------------------------------------------------------
const (
	RXST_STX    = 0
	RXST_LEN    = RXST_STX + 1
	RXST_SEQ    = RXST_LEN + 1
	RXST_OPCODE = RXST_SEQ + 1
	RXST_DATA   = RXST_OPCODE + 1
)

// ---------------------------------------------------------------------------
// Constant
// ---------------------------------------------------------------------------
const (
	MAX_BUFLEN       = 10240
	SHORTBUF_LEN     = 32
	CONNECT_INTERVAL = 2000 // 2sec	// 커넥트 요청 주기

	CHBUF_CommToPrc  = 300
	CHBUF_HttpToComm = 300
)

const (
	ComRecvBuf    = 1
	ComSendBuf    = 100
	CHBUF_NET_PRC = 300
)

// ---------------------------------------------------------------------------
// Recevied Network message for sending msgapp
// ---------------------------------------------------------------------------
type NetworkMSG struct {
	Id     int
	Code   byte
	Data   []byte    // 통신데이터
	Evt    []byte    // 서버에서 발생시킨 이벤트 데이터
	RecvTm time.Time // 수신 시간
}

// ---------------------------------------------------------------------------
// Http request and command message
// ---------------------------------------------------------------------------
type CommReqMSG struct {
	Key    int   // session key
	Code   byte  // OPCODE
	Ids    []int // id list
	Data   []byte
	Rst    int       // 결과	0:전송안함, 1:응답 대기중 2:ACK, >3:NACK
	SendTm time.Time // record send time
}

type MessageQue struct {
	Code      byte
	M_message []byte
	Seq       byte
}

// ---------------------------------------------------------------------------
// CommHandler
// ---------------------------------------------------------------------------
type CommHandler struct {
	Id  int
	Tcp comm.TcpHandler

	isCommenvSet bool

	msg_q *dt.SlQueye[MessageQue]

	m_message    []byte
	m_length     int
	connectTimer time.Time
	stateTimer   time.Time

	m_index   int
	m_rxState int
	m_txSeq   byte
	m_rxSeq   byte
}

// ------------------------------------------------------------------------------
// MakeNetHandler
// ------------------------------------------------------------------------------
func MakeNetHandler(name string, port int, addr string) *CommHandler {

	a := &CommHandler{
		Tcp:       comm.NewTcpHandler(name, port, addr),
		m_message: []byte{},
		m_length:  0,
	}

	return a
}
