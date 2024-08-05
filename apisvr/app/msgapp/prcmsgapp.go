package msgapp

import (
	"apisvr/app/am"
	"apisvr/app/netapp"
	"time"

	"github.com/gdygd/goglib"
)

// ------------------------------------------------------------------------------
// PrcmsgAppHandler
// ------------------------------------------------------------------------------
type PrcmsgAppHandler struct {
}

// ------------------------------------------------------------------------------
// processTest
// ------------------------------------------------------------------------------
func (a *PrcmsgAppHandler) processTest(data netapp.NetworkMSG) {
	am.Applog.Print(2, "processTest..id:%d, code:[%02X], data:[%v]", data.Id, data.Code, data.Data)

}

// ------------------------------------------------------------------------------
// msgHandler
// ------------------------------------------------------------------------------
func (a *PrcmsgAppHandler) msgHandler(data netapp.NetworkMSG) {
	// process comm packet

	switch data.Code {

	case netapp.PT_TEST:
		// vms controller state
		a.processTest(data)

	default:
		am.Applog.Warn("[prcapp] prcmsgapp Undefined opcode (%d)%02x", data.Id, data.Code)
	}

}

// ------------------------------------------------------------------------------
// THRPrcmsg
// ------------------------------------------------------------------------------
func THRPrcmsg(t *goglib.Thread, chThrStop chan bool, arg1, arg2, arg3 interface{}) {
	am.Applog.Always("THRPrcmsg start")

	var terminate = false
	app := &PrcmsgAppHandler{}

	//------------------------------------
	// tx routine
	//------------------------------------
	var runcnt int = 0
	for {
		select {
		case <-chThrStop:
			am.Applog.Always("[1]thread stop,, process message thread quit..")
			terminate = true
			break
		default:
			// Pop Network message
			netmsg, ok := netapp.PopNetMsg()
			if ok {
				app.msgHandler(netmsg)
			}

		}

		if terminate {
			break
		}

		if runcnt%1000 == 0 {
			am.Applog.Print(2, "message process[%d]", runcnt)
		}
		runcnt++

		t.RunBase.UpdateRunInfo()
		time.Sleep(time.Millisecond * 100)
	}

	am.Applog.Always("[2]message process thread quit..")
}
