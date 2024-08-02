package httpapp

import (
	"apisvr/app/am"
	"apisvr/app/netapp"
	"sync"
	"time"

	"github.com/gdygd/goglib"
)

// ---------------------------------------------------------------------------
// Mutex
// ---------------------------------------------------------------------------
var rwMCtrlSession = new(sync.RWMutex) // read, write mutex

// ---------------------------------------------------------------------------
// Ctrl command session
// ---------------------------------------------------------------------------
var ctrlSessionKeyList []SessionObj = makeCtrlSessionKey()

func isTcpChanClosed(ch <-chan netapp.CommReqMSG) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func MakeCtrlMessage(sessionKey, userIdx int, code byte, Ids []int, datas []int) netapp.CommReqMSG {
	var msg netapp.CommReqMSG = netapp.CommReqMSG{}

	msg.Key = sessionKey

	msg.Ids = make([]int, 0)
	msg.Data = make([]byte, 0)
	msg.Code = byte(code)

	for _, id := range Ids {
		msg.Ids = append(msg.Ids, id)
	}

	for _, dt := range datas {
		msg.Data = append(msg.Data, byte(dt))
	}

	return msg
}

func makeCtrlSessionKey() []SessionObj {
	sessionList := make([]SessionObj, 0)
	for idx := 1; idx <= 1000; idx++ {
		var session = SessionObj{Key: idx}

		sessionList = append(sessionList, session)
	}

	return sessionList
}

func GetSessionKey() int {

	rwMCtrlSession.Lock()

	var frontKeyObj SessionObj = SessionObj{}
	if len(ctrlSessionKeyList) > 0 {

		frontKeyObj, ctrlSessionKeyList = ctrlSessionKeyList[0], ctrlSessionKeyList[1:]
		am.Applog.Print(5, "GetCtrlSessionKey:%v", frontKeyObj)
	}
	rwMCtrlSession.Unlock()

	// 세션 존재 여부 확인, close일 경우 channel make
	if isTcpChanClosed(netapp.ChToHttp[frontKeyObj.Key]) {
		netapp.ChToHttp[frontKeyObj.Key] = make(chan netapp.CommReqMSG, 1)
	}
	return frontKeyObj.Key
}

func ClearSessionKey(key int) {
	rwMCtrlSession.Lock()
	ctrlSessionKeyList = append(ctrlSessionKeyList, SessionObj{Key: key})
	rwMCtrlSession.Unlock()
}

func WaitResult(sessionKey int) netapp.CommReqMSG {

	var resultData netapp.CommReqMSG = netapp.CommReqMSG{}
	var rstTimer time.Time = time.Time{}
	rstTimer = time.Now()

WaitLoop:
	for {
		select {
		case rstData := <-netapp.ChToHttp[sessionKey]:
			resultData = rstData
			am.Applog.Print(6, "Ctrl command rst Data : %v", rstData)
			break WaitLoop
		default:
			// max timeout
			if goglib.CheckElapsedTime(&rstTimer, 5000) {

				// 타임 아웃시 채널 close
				if isTcpChanClosed(netapp.ChToHttp[sessionKey]) {
					close(netapp.ChToHttp[sessionKey])
				}

				resultData.Rst = BE_TIMEOUT
				am.Applog.Print(6, "[BE_TIMEOUT]Ctrl command rst Data : %v", resultData.Rst)
				//result = BE_TIMEOUT
				break WaitLoop
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	return resultData
}
