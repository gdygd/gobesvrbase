package netapp

import (
	"apisvr/app/am"
	"time"

	"github.com/gdygd/goglib"
)

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initClientEnv(client *CommHandler) bool {
	am.Netlog.Always("initClientEnv")

	client.connectTimer = time.Now()
	client.stateTimer = time.Now()

	return true
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func clearClientEnv(client *CommHandler) {
	am.Netlog.Always("clearClientEnv")
}

// ------------------------------------------------------------------------------
// THRNetclient
// ------------------------------------------------------------------------------
func THRNetclient(t *goglib.Thread, chThrStop chan bool, arg1, arg2, arg3 interface{}) {
	am.Netlog.Print(2, "THRNetclient start...")

	am.Netlog.Always("Start THRNetclient..")
	a, ok := arg1.(*CommHandler)

	if !ok {
		am.Applog.Error("THRNetclient arg error... ")
		return
	}

	chquit := make(chan bool, 1)
	var terminate = false
	var initOk bool = false
	initOk = initClientEnv(a)
	defer clearClientEnv(a)

	//------------------------------------
	// tx routine
	//------------------------------------
	go func() {
		for {
			select {
			case <-chThrStop:
				am.Netlog.Always("[1]thread stop,, Client tx thread quit..")
				terminate = true
				chquit <- true
				break
			default:
				//
			}

			if terminate {
				break
			}

			t.RunBase.UpdateRunInfo()
			a.Manage()
			if a.Tcp.IsConnected() {
				a.ManageTx()
				t.RunBase.UpdateRunInfo()
			}
			time.Sleep(time.Millisecond * 100)
		}

		am.Netlog.Always("[2]thread stop,, Client rx thread quit..")
	}()

	//------------------------------------
	// Rx routine & run manage
	//------------------------------------
	for {
		if !initOk {
			am.Netlog.Warn("Client initEnv fail...")
			break
		}

		select {
		case <-chquit:
			am.Netlog.Always("Client rx thread quit..")
			break
		default:
			//
		}

		if terminate {
			am.Netlog.Always("Client rx thread quit..")
			break
		}

		t.RunBase.UpdateRunInfo()
		a.Manage()
		if a.Tcp.IsConnected() {
			a.ManageRx()
		}
		time.Sleep(time.Millisecond * 100)
	}

}
