package kafkaapp

import (
	"apisvr/app/am"
	"sync"
)

// ---------------------------------------------------------------------------
// mutex
// ---------------------------------------------------------------------------
var msgMTX = &sync.Mutex{} // comm manage mutex

// kafka message -> procMsg 채널
var CHkafkamsg chan KafkakMSG = initChannelKafkaMsg() // kafka <-> procMsg 채널

// ---------------------------------------------------------------------------
// initChannelNetMsg
// ---------------------------------------------------------------------------
func initChannelKafkaMsg() chan KafkakMSG {
	channel := make(chan KafkakMSG, CHBUF_NET_PRC)
	return channel
}

func clearKafkaMsgChannel() {
	for len(CHkafkamsg) > 0 {
		//<-CHkafkamsg
		PopKafkaMsg()
	}
}

func checkKafkaMsgChannel() {
	if len(CHkafkamsg) >= CHBUF_NET_PRC {
		am.Netlog.Warn("CHkafkamsg clear.. [%d][%d]", len(CHkafkamsg), CHBUF_NET_PRC)
		clearKafkaMsgChannel()
	}
}

func SendKafkaToMsgProc(msg KafkakMSG) bool {
	//am.Netlog.Print(3, "SendKafkaToMsgProc(1)")
	msgMTX.Lock()
	checkKafkaMsgChannel()
	CHkafkamsg <- msg
	msgMTX.Unlock()

	//am.Netlog.Print(3, "SendKafkaToMsgProc(2))")

	return true
}

func PopKafkaMsg() (KafkakMSG, bool) {
	var msg KafkakMSG
	var ok bool = false
	msgMTX.Lock()
	//am.Netlog.Print(3, "PopKafkaMsg...")
	if len(CHkafkamsg) > 0 {
		//am.Netlog.Print(3, "PopKafkaMsg... Pop(1)")
		msg = <-CHkafkamsg
		//msgMTX.Unlock()
		ok = true
		//return msg, true
	}

	msgMTX.Unlock()
	//am.Netlog.Print(3, "PopKafkaMsg... Pop(2)[%v]", ok)
	return msg, ok
}
