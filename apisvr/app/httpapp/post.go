package httpapp

import (
	"apisvr/app/am"
	"apisvr/app/netapp"
	"encoding/json"
	"net/http"
)

// ------------------------------------------------------------------------------
// PostTest, TEST
// ------------------------------------------------------------------------------
func (a *HttpAppHandler) PostTest(w http.ResponseWriter, r *http.Request) {
	var returnData ResponseMsg = ResponseMsg{}
	var result int = SUCCESS

	am.Applog.Print(2, "PostTest")
	returnData.Data = "PostTest"

	returnData.Result = result
	json.NewEncoder(w).Encode(returnData)

}

// ------------------------------------------------------------------------------
// NetCmdTest, TEST
// ------------------------------------------------------------------------------
func (a *HttpAppHandler) NetCmdTest(w http.ResponseWriter, r *http.Request) {
	var returnData ResponseMsg = ResponseMsg{}
	var result int = SUCCESS

	am.Applog.Print(2, "NetCmdTest")
	returnData.Data = "NetCmdTest"

	// Get http session key
	sessionKey := GetSessionKey()
	defer ClearSessionKey(sessionKey)

	var packetBuf []int = []int{}
	// Make mesage info (test)
	msg := MakeCtrlMessage(sessionKey, 0, byte(netapp.PT_TEST), []int{0}, packetBuf)

	ok := netapp.SendHttpToNet(msg)

	// waith for until netapp response
	var resultData netapp.CommReqMSG
	if ok {
		resultData = WaitResult(sessionKey)
	}

	returnData.Result = result
	returnData.Data = resultData
	json.NewEncoder(w).Encode(returnData)
}
