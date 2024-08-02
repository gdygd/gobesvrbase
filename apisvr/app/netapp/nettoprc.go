package netapp

import (
	"apisvr/app/am"
	"sync"
)

// ---------------------------------------------------------------------------
// mutex
// ---------------------------------------------------------------------------
var netmsgMTX = &sync.Mutex{} // comm manage mutex

// tcp message -> procMsg 채널
var CHnetmsg chan NetworkMSG = initChannelNetMsg() // tcp <-> procMsg 채널

// ---------------------------------------------------------------------------
// initChannelNetMsg
// ---------------------------------------------------------------------------
func initChannelNetMsg() chan NetworkMSG {
	channel := make(chan NetworkMSG, CHBUF_NET_PRC)
	return channel
}

func clearNetMsgChannel() {
	for len(CHnetmsg) > 0 {
		//<-CHnetmsg
		PopNetMsg()
	}
}

func checkNetMsgChannel() {
	if len(CHnetmsg) >= CHBUF_NET_PRC {
		am.Netlog.Warn("CHnetmsg clear.. [%d][%d]", len(CHnetmsg), CHBUF_NET_PRC)
		clearNetMsgChannel()
	}
}

func SendNetMsgToMsgProc(msg NetworkMSG) bool {
	//am.Netlog.Print(3, "SendNetMsgToMsgProc(1)")
	netmsgMTX.Lock()
	checkNetMsgChannel()
	CHnetmsg <- msg
	netmsgMTX.Unlock()

	//am.Netlog.Print(3, "SendNetMsgToMsgProc(2))")

	return true
}

func PopNetMsg() (NetworkMSG, bool) {
	var msg NetworkMSG
	var ok bool = false
	netmsgMTX.Lock()
	//am.Netlog.Print(3, "PopNetMsg...")
	if len(CHnetmsg) > 0 {
		//am.Netlog.Print(3, "PopNetMsg... Pop(1)")
		msg = <-CHnetmsg
		//netmsgMTX.Unlock()
		ok = true
		//return msg, true
	}

	netmsgMTX.Unlock()
	//am.Netlog.Print(3, "PopNetMsg... Pop(2)[%v]", ok)
	return msg, ok
}
