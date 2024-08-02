package netapp

// http <-> net client간 제어 메시지 관리 채널

// key : object id
var ChFromHttp chan CommReqMSG = initChannelFromHttp() // http -> tcp send

// key : session key
var ChToHttp map[int]chan CommReqMSG = initChannelToHttp() // http <- tcp recv

// ---------------------------------------------------------------------------
// initChannelFromHttp
// ---------------------------------------------------------------------------
func initChannelFromHttp() chan CommReqMSG {
	channel := make(chan CommReqMSG, CHBUF_HttpToComm)

	return channel
}

// ---------------------------------------------------------------------------
// initChannelToHttp
// ---------------------------------------------------------------------------
func initChannelToHttp() map[int]chan CommReqMSG {
	mpChannel := make(map[int]chan CommReqMSG)
	//세션 key별 채널 생성
	for key := 1; key <= 1000; key++ {
		mpChannel[key] = make(chan CommReqMSG, 1)
	}

	return mpChannel
}

func clearToHttpChan(key int) {
	for len(ChToHttp[key]) > 0 {
		<-ChToHttp[key]
	}
}

func CheckToHttpChan(key int) {
	if len(ChToHttp[key]) >= ComRecvBuf {
		clearToHttpChan(key)
	}
}

func clearFromHttpChan() {
	for len(ChFromHttp) > 0 {
		<-ChFromHttp
	}
}

func CheckFromHttpChan() {
	if len(ChFromHttp) >= ComSendBuf {
		clearFromHttpChan()
	}
}

func SendHttpToNet(msg CommReqMSG) bool {
	// http -> client

	CheckFromHttpChan()
	ChFromHttp <- msg

	return true
}

func SendNetToHttp(msg CommReqMSG, key int) bool {
	// http <- client

	CheckToHttpChan(key)
	ChToHttp[key] <- msg

	return true
}
