package netapp

import (
	"apisvr/app/am"
	"sync"
	"time"

	"github.com/gdygd/goglib"
	dt "github.com/gdygd/goglib/datastructure"
)

// ---------------------------------------------------------------------------
// Mutex
// ---------------------------------------------------------------------------
var mngMTX = &sync.Mutex{} // comm manage mutex

// ------------------------------------------------------------------------------
// SetTcpCommEnv
// ------------------------------------------------------------------------------
func (c *CommHandler) SetTcpCommEnv(addr string, port int) {

	am.Netlog.Print(2, "SetTcpCommEnv : %s, %d", addr, port)
	c.Tcp.SetCommEnv("Client", port, addr)

}

// ------------------------------------------------------------------------------
// GetTcpCommEnv
// ------------------------------------------------------------------------------
func (c *CommHandler) GetTcpCommEnv() (string, int) {
	addr := c.Tcp.GetAddress()
	port := c.Tcp.GetPort()

	return addr, port

}

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (c *CommHandler) Open() bool {

	am.Netlog.Print(1, "Open...")

	if !goglib.CheckElapsedTime(&c.connectTimer, CONNECT_INTERVAL) {
		return true
	}

	// manage connect
	if c.Tcp.IsConnected() {
		am.Netlog.Print(2, "Open..IsConnected")
		return true
	}

	am.Netlog.Print(2, "is not open.. try open tcp..")
	am.Netlog.Print(2, "Network Connect...")

	ok, err := c.Tcp.Connect()
	if ok && err == nil {
		am.Netlog.Print(3, "Network connected!")
		c.ResetCommEnv()
	} else {
		am.Netlog.Print(2, "Connect Error: %v, %v", ok, err)
		ok = false
	}

	return ok
}

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (c *CommHandler) OpenImmediate() bool {

	am.Netlog.Print(2, "OpenImmediate...")

	// if !goglib.CheckElapsedTime(&c.connectTimer, CONNECT_INTERVAL) {
	// 	return true
	// }

	// manage connect
	if c.Tcp.IsConnected() {
		am.Netlog.Print(2, "Open..IsConnected")
		return true
	}

	am.Netlog.Print(2, "is not open.. try open tcp..")
	am.Netlog.Print(2, "Network Connect...")

	ok, err := c.Tcp.Connect()
	if ok && err == nil {
		am.Netlog.Print(3, "Network connected!")
		c.ResetCommEnv()
	} else {
		am.Netlog.Print(2, "Connect Error: %v, %v", ok, err)
		ok = false
	}

	return ok
}

// ---------------------------------------------------------------------------
// ResetCommEnv
// ---------------------------------------------------------------------------
func (c *CommHandler) ResetCommEnv() {

	c.m_message = c.m_message[:0]
	c.m_index = 0
	c.m_length = 0
	c.m_rxState = RXST_STX

	c.m_txSeq = 0
	c.m_rxSeq = 0

	c.msg_q = dt.NewSlQueue[MessageQue]()
}

// ---------------------------------------------------------------------------
// ResetRxStatus
// ---------------------------------------------------------------------------
func (c *CommHandler) ResetRxStatus() {
	// clear
	c.m_message = c.m_message[:0]
	c.m_index = 0
	c.m_rxState = RXST_STX
}

// ---------------------------------------------------------------------------
// SetRxStatus
// ---------------------------------------------------------------------------
func (c *CommHandler) SetRxStatus(rxState int, readCount int) {
	c.m_rxState = rxState
	c.m_index += readCount

	if rxState == RXST_STX {
		c.ResetRxStatus()
	}
}

// ------------------------------------------------------------------------------
// Manage
// ------------------------------------------------------------------------------
func (c *CommHandler) Manage() {

	mngMTX.Lock()
	// is not Open
	//	Open()
	if !c.IsOpen() {
		am.Netlog.Print(1, "is not open.. try open tcp..")
		c.ResetCommEnv()
		c.Open()
	}

	mngMTX.Unlock()

}

// ------------------------------------------------------------------------------
// ManageRx
// ------------------------------------------------------------------------------
func (c *CommHandler) ManageRx() bool {
	am.Netlog.Print(1, "ManageRX..")
	rxbuf := make([]byte, MAX_BUFLEN, MAX_BUFLEN)

	if c.Tcp.Connected {
		n, err := c.Tcp.Read(rxbuf)

		if err != nil {
			am.Netlog.Error("Rx Error [%d],%v", n, err.Error())
			c.Close()
		} else {
			am.Netlog.Dump(2, "RXLow", rxbuf[:n], n)
			//c.RxHandler(rxbuf[:n], n)
		}
	}

	am.Netlog.Print(2, "ManageRX..end..")
	return true
}

// ------------------------------------------------------------------------------
// RxHandler
// ------------------------------------------------------------------------------
func (c *CommHandler) RxHandler(data []byte, length int) {

	if length >= MAX_BUFLEN {
		am.Netlog.Error("RX lenght over max size..(maxsize : %d) : len:%d", MAX_BUFLEN, length)
		return
	}

	m_msg := &c.m_message
	for index := 0; index < length; index++ {
		m_rxstt := c.m_rxState

		if m_rxstt == RXST_STX {
			am.Netlog.Print(1, "STX.. index:%d", index)
			if data[index] != PT_STX {
				am.Netlog.Warn("Invalid STX..%02X", data[index])
				c.SetRxStatus(RXST_STX, 0)
				break
			}
			am.Netlog.Print(1, "RX STX")
			*m_msg = append(*m_msg, data[index])
			c.SetRxStatus(m_rxstt+1, 1)

		} else if m_rxstt == RXST_LEN {
			am.Netlog.Print(1, "LEN.. index:%d", index)
			am.Netlog.Print(1, "RX LEN")
			if len(data) > index+2 {
				*m_msg = append(*m_msg, data[index:index+2]...)
				c.SetRxStatus(m_rxstt+1, 2)
				index += 1
			} else {
				am.Netlog.Print(2, "len size invalid..")
				continue
			}

		} else if m_rxstt == RXST_SEQ {
			am.Netlog.Print(1, "SEQ.. index:%d", index)
			am.Netlog.Print(1, "RX SEQ")
			*m_msg = append(*m_msg, data[index])
			c.SetRxStatus(m_rxstt+1, 1)

		} else if m_rxstt == RXST_OPCODE {
			am.Netlog.Print(1, "OPCODE.. index:%d", index)
			am.Netlog.Print(1, "RX CODE")
			*m_msg = append(*m_msg, data[index])
			c.SetRxStatus(m_rxstt+1, 1)

		} else if m_rxstt == RXST_DATA {
			am.Netlog.Print(1, "DATA.. index:%d", index)
			am.Netlog.Print(1, "RX DATA")
			am.Netlog.Dump(1, "RXDATA", *m_msg, len(*m_msg))
			c.m_length = goglib.GetNumber(*m_msg, PT_LEN_IDX, 2, goglib.ED_BIG)
			needsize := c.m_length - c.m_index
			remainsize := length - index // 현재 패킷의 남은 부분

			am.Netlog.Print(1, "needsize : %d, reaminpacket size : %d", needsize, remainsize)
			if needsize <= remainsize {
				*m_msg = append(*m_msg, data[index:index+needsize]...)
				index += needsize - 1
				c.SetRxStatus(RXST_DATA, needsize)

			} else { //
				// needsize > remainsize
				*m_msg = append(*m_msg, data[index:index+remainsize]...)
				index += remainsize - 1
				c.SetRxStatus(RXST_DATA, remainsize)
			}

			if c.m_length == c.m_index {
				// check crc

				// msg handler
				c.MsgHandler()

				c.SetRxStatus(RXST_STX, 0)
			}
		}
	}
}

// ------------------------------------------------------------------------------
// MsgHandler
// ------------------------------------------------------------------------------
func (c *CommHandler) MsgHandler() {
	code := c.m_message[PT_OP_IDX]
	am.Netlog.Print(1, "recv opcode %02X", code)

	var msg []byte = make([]byte, len(c.m_message))
	copy(msg, c.m_message)
	switch code {
	case PT_TEST:
		c.ProcessTest(msg)
	default:
		am.Netlog.Warn("Invalid opcode received..[%02X]", code)
	}
}

// ---------------------------------------------------------------------------
// ProcessLocStt2
// ---------------------------------------------------------------------------
func (c *CommHandler) ProcessTest(message []byte) {
	// test protocol

	locid := goglib.GetNumber(message, PT_DATA_IDX, 4, goglib.ED_BIG)
	seq := message[PT_SEQ_IDX]

	am.Netlog.Print(2, "ProcessTest..%d, %d", locid, seq)

	code := message[PT_OP_IDX]

	var msg NetworkMSG = NetworkMSG{}
	msg.Code = byte(code)
	msg.Data = make([]byte, len(message))
	copy(msg.Data, message[:(len(message))])
	am.Netlog.Print(1, "ProcessLocStt chLen(%d) [%02X]", len(CHnetmsg), msg.Code)

	SendNetMsgToMsgProc(msg)

	//c.SendAck(locid, seq)
}

// ------------------------------------------------------------------------------
// ManageTx
// ------------------------------------------------------------------------------
func (c *CommHandler) ManageTx() bool {
	//rxbuf := make([]byte, MAX_BUFLEN, MAX_BUFLEN)

	var ok bool = true

	if !c.Tcp.Connected {
		ok = false
		return ok
	}

	if goglib.CheckElapsedTime(&c.stateTimer, 1000) {
		c.SendTest()
	}

	msg, isexist := c.msg_q.Pop()
	if isexist {
		c.SendMsg(msg)
	}

	select {
	case ctrlMsg := <-ChFromHttp:
		//am.Netlog.Print(7, "ManageTx(2)")
		ok = c.TxCtrlHandler(ctrlMsg)
	default:
		//am.Netlog.Print(7, "ManageTx(3)")
		ok = true
	}

	return true
}

// ------------------------------------------------------------------------------
// SendAck
// ------------------------------------------------------------------------------
func (c *CommHandler) SendTest() {
	// send only time_t
	var info []byte = make([]byte, MAX_BUFLEN)
	var idx int = 0

	curtm := time.Now()
	utcsecs := curtm.Unix()

	// time
	goglib.SetNumber64(info, idx, utcsecs, 8, goglib.ED_BIG)
	idx += 8

	c.PushTxData(PT_TEST, info[:idx])
}

// ------------------------------------------------------------------------------
// SendAck
// ------------------------------------------------------------------------------
func (c *CommHandler) SendAck(locid int, seq byte) {
	var info []byte = make([]byte, SHORTBUF_LEN)
	var idx int = 0

	var msg MessageQue = MessageQue{
		Code:      PT_ACK,
		M_message: info[:idx],
		Seq:       seq,
	}

	c.msg_q.Push(msg)
}

// ------------------------------------------------------------------------------
// PushTxData
// ------------------------------------------------------------------------------
func (c *CommHandler) PushTxData(code byte, buf []byte) {
	am.Netlog.Print(2, "PushTxData ..")

	if c.m_txSeq >= 255 {
		c.m_txSeq = 1
	}

	var msg MessageQue = MessageQue{
		Code:      code,
		M_message: buf,
		Seq:       c.m_txSeq,
	}
	c.m_txSeq++

	c.msg_q.Push(msg)
}

// ------------------------------------------------------------------------------
// TxCtrlHandler
// ------------------------------------------------------------------------------
func (c *CommHandler) TxCtrlHandler(data CommReqMSG) bool {

	ctrlCode := data.Code
	am.Netlog.Print(3, "[TxCtrlHandler] (%02X) [%v]", ctrlCode, data.Data)
	ok := true

	switch ctrlCode {
	case PT_TEST:
		am.Netlog.Print(2, "TxCtrlHandler TSPPC_TX_TEST..")
		c.PushTxData(data.Code, data.Data)

		//test net to msgproc
		var testpacket []byte = []byte{2, 0, 11, 7, 153, 1, 2, 3, 4, 5, 3}
		c.ProcessTest(testpacket)

	default:
		am.Netlog.Warn("Undefined ctrl message code [%02x]", ctrlCode)
	}

	return ok
}

// ------------------------------------------------------------------------------
// SendMsg
// ------------------------------------------------------------------------------
func (c *CommHandler) SendMsg(msg MessageQue) bool {
	len := len(msg.M_message)
	return (c.SendMessage(msg.Code, msg.Seq, len, msg.M_message, 3))
}

// ---------------------------------------------------------------------------
// SendMessage
// ---------------------------------------------------------------------------
func (c *CommHandler) SendMessage(code, seq byte, length int, info []byte, lv int) bool {

	//-----------------------------------------------
	// General packet
	//-----------------------------------------------
	minPacket := HEAD_SIZE + length + TAIL_SIZE

	packet := make([]byte, minPacket, minPacket)
	hPacket := make([]byte, HEAD_SIZE, HEAD_SIZE)

	packetlen := HEAD_SIZE + length + TAIL_SIZE

	hPacket[PT_STX_IDX] = PT_STX                                       // stx
	goglib.SetNumber(hPacket, PT_LEN_IDX, packetlen, 2, goglib.ED_BIG) //len
	hPacket[PT_SEQ_IDX] = seq
	hPacket[PT_OP_IDX] = code

	// set header data
	copy(packet, hPacket)

	// set data packet
	for i := 0; i < length; i++ {
		packet[PT_DATA_IDX+i] = info[i]
	}

	am.Netlog.Print(2, "Tx size:%d min size:%d", len(packet), minPacket)
	am.Netlog.Dump(lv, "TX raw#1", packet, len(packet))

	packet[minPacket-1] = PT_ETX

	am.Netlog.Dump(2, "TX raw#1", packet, len(packet))

	_, err := c.Tcp.Send(packet)
	if err != nil {
		am.Netlog.Error("send err..%v", err)
	}

	return true
}

// ------------------------------------------------------------------------------
// IsOpen
// ------------------------------------------------------------------------------
func (c *CommHandler) IsOpen() bool {
	//mngMTX.Lock()

	isopen := c.Tcp.IsConnected()
	//mngMTX.Unlock()
	return isopen
}

// ------------------------------------------------------------------------------
// Close
// ------------------------------------------------------------------------------
func (c *CommHandler) Close() {
	am.Netlog.Print(1, "Tcp close...")

	if !c.IsOpen() {
		return
	}
	c.msg_q.QClear()
	c.Tcp.Close()
	am.Netlog.Print(2, "Tcp close...ok")
}
