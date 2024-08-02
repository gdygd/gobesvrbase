package netapp

import (
	"apisvr/app/am"

	"github.com/gdygd/goglib"
)

// ------------------------------------------------------------------------------
// THRVmsserver
// ------------------------------------------------------------------------------
func THRNetclient(t *goglib.Thread, chThrStop chan bool, arg1, arg2, arg3 interface{}) {
	am.Netlog.Print(2, "THRNetclient start...")

}
