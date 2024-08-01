package httpapp

import (
	"apisvr/app/am"
	"sync"

	"github.com/gdygd/goglib"
)

// ---------------------------------------------------------------------------
// Mutex
// ---------------------------------------------------------------------------
var rwMSSeSession = new(sync.RWMutex) // read, write vds session mutex
// ---------------------------------------------------------------------------
// Mutex(SSE mutex)
// ---------------------------------------------------------------------------
var sseMsgMTX = &sync.Mutex{} // sse message mutex

// ---------------------------------------------------------------------------
// SSE session key list
// ---------------------------------------------------------------------------
var sseSessionKeyList []SessionObj = makeSSeSessionKey()
var ActivesseSessionList []SessionObj = []SessionObj{}

// ---------------------------------------------------------------------------
// SSE channel
// ---------------------------------------------------------------------------
var SseMsgChan map[int]chan goglib.EventData = initSSeMsgChannel()

// ---------------------------------------------------------------------------
// makeSSeSessionKey
// ---------------------------------------------------------------------------
func makeSSeSessionKey() []SessionObj {
	sseKeyList := make([]SessionObj, 0)

	for idx := 1; idx <= goglib.CHBUF_SSE; idx++ {
		var session = SessionObj{Key: idx}
		sseKeyList = append(sseKeyList, session)
	}

	return sseKeyList
}

// ---------------------------------------------------------------------------
// initSSeMsgChannel
// ---------------------------------------------------------------------------
func initSSeMsgChannel() map[int]chan goglib.EventData {
	mpChannel := make(map[int]chan goglib.EventData)
	//세션 key별 채널 생성
	for key := 1; key <= goglib.CHBUF_SSE; key++ {
		mpChannel[key] = make(chan goglib.EventData, goglib.CHBUF_SSE)
	}

	return mpChannel
}

// ---------------------------------------------------------------------------
// PopSSEMsgChannel
// ---------------------------------------------------------------------------
func PopSSEMsgChannel(key int) (goglib.EventData, bool) {
	var popData goglib.EventData
	sseMsgMTX.Lock()
	if len(SseMsgChan[key]) > 0 {
		popData = <-SseMsgChan[key]
		sseMsgMTX.Unlock()

		return popData, true
	}
	sseMsgMTX.Unlock()

	return popData, false
}

func clearSSEMsgChannel(key int) {
	for len(SseMsgChan[key]) > 0 {
		PopSSEMsgChannel(key)
	}
}

func CheckSSEMsgChannel(key int) {
	if len(SseMsgChan[key]) >= goglib.CHBUF_SSE {
		clearSSEMsgChannel(key)
	}
}

func removeActiveSSeKey(index int) {
	ActivesseSessionList = append(ActivesseSessionList[:index], ActivesseSessionList[index+1:]...)
}

// ---------------------------------------------------------------------------
// GetSSeSessionKey
// ---------------------------------------------------------------------------
func GetSSeSessionKey() int {
	rwMSSeSession.Lock()

	var frontKeyObj SessionObj = SessionObj{}

	if len(sseSessionKeyList) > 0 {

		frontKeyObj, sseSessionKeyList = sseSessionKeyList[0], sseSessionKeyList[1:]
		ActivesseSessionList = append(ActivesseSessionList, frontKeyObj)

		am.Applog.Print(2, "GetSSeSessionKey:%v", frontKeyObj)
	}

	rwMSSeSession.Unlock()

	return frontKeyObj.Key
}

// ---------------------------------------------------------------------------
// ClearSSeSessionKey
// ---------------------------------------------------------------------------
func ClearSSeSessionKey(key int) {
	rwMSSeSession.Lock()

	sseSessionKeyList = append(sseSessionKeyList, SessionObj{Key: key})
	am.Applog.Print(2, "ClearSSeSessionKeySeSS:[%v]", sseSessionKeyList)
	am.Applog.Print(2, "ClearSSeSessionKeyACT1:[%d][%v]", key, ActivesseSessionList)
	for index, data := range ActivesseSessionList {
		if data.Key == key {
			removeActiveSSeKey(index)
			break
		}
	}

	am.Applog.Print(2, "ClearSSeSessionKeyACT2:[%d][%v]", key, ActivesseSessionList)

	rwMSSeSession.Unlock()
}
