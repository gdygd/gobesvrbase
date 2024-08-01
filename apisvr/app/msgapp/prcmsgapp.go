package msgapp

import (
	"apisvr/app/am"
	"time"

	"github.com/gdygd/goglib"
)

// ------------------------------------------------------------------------------
// PrcmsgAppHandler
// ------------------------------------------------------------------------------
type PrcmsgAppHandler struct {
}

// ------------------------------------------------------------------------------
// processtest
// ------------------------------------------------------------------------------
func (a *PrcmsgAppHandler) processtest() {

}

// ------------------------------------------------------------------------------
// msgHandler
// ------------------------------------------------------------------------------
func (a *PrcmsgAppHandler) msgHandler() {

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
	for {
		select {
		case <-chThrStop:
			am.Applog.Always("[1]thread stop,, process message thread quit..")
			terminate = true
			break
		default:
			//
			app.msgHandler()
		}

		if terminate {
			break
		}

		t.RunBase.UpdateRunInfo()
		time.Sleep(time.Millisecond * 50)
	}

	am.Applog.Always("[2]message process thread quit..")
}
