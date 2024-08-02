package httpapp

import (
	"apisvr/app/am"
	"encoding/json"
	"net/http"
)

// ------------------------------------------------------------------------------
// GetTest, get
// ------------------------------------------------------------------------------
func (a *HttpAppHandler) GetTest(w http.ResponseWriter, r *http.Request) {

	var returnData ResponseMsg = ResponseMsg{}
	var result int = SUCCESS

	am.Applog.Print(2, "GetTest")

	readInfo, err := a.dbHnd.ReadTest()

	// check error
	if err != nil {
		am.Applog.Error("GetLocalInfo, /local error:%v", err)
		result = FAIL
		returnData.Result = result
		json.NewEncoder(w).Encode(returnData)
		return
	}

	returnData.Data = readInfo
	returnData.Result = result
	json.NewEncoder(w).Encode(returnData)
}
