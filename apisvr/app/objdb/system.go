package objdb

import (
	"apisvr/app/am"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gdygd/goglib"
)

type SystemInfo struct {
	svrTime    time.Time `json:"-"`
	StrSvrtime string    `json:"svrtime"`
}

// ---------------------------------------------------------------------------
// Global
// ---------------------------------------------------------------------------
var SysInfo *SystemInfo = &SystemInfo{svrTime: time.Now()}

// ------------------------------------------------------------------------------
//
//	SetSvrTime
//
// ------------------------------------------------------------------------------
func (s *SystemInfo) SetSvrTime(utc int64) {
	s.svrTime = time.Unix(utc, 0)
	s.StrSvrtime = fmt.Sprintf("%04d-%02d-%02d  %02d:%02d:%02d", s.svrTime.Year(), s.svrTime.Month(), s.svrTime.Day(), s.svrTime.Hour(), s.svrTime.Minute(), s.svrTime.Second())

	s.sseSystemInfo()
}

func (s *SystemInfo) GetSvrTime() string {
	return s.StrSvrtime
}

// ---------------------------------------------------------------------------
// SSESystemInfo
// ---------------------------------------------------------------------------
func (s *SystemInfo) sseSystemInfo() {

	svrtime := s.GetSvrTime()

	b, _ := json.Marshal(svrtime)
	var evdata goglib.EventData = goglib.EventData{}
	evdata.Msgtype = EVCD_SVRTIME
	evdata.Data = string(b)
	evdata.Id = "1"

	am.Applog.Print(2, "sseSystemInfo [%s](%v)", evdata.Msgtype, evdata.Data)
	goglib.SendSSE(evdata)

}
