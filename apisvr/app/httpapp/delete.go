package httpapp

import (
	"apisvr/app/am"
	"encoding/json"
	"net/http"
)

// ------------------------------------------------------------------------------
// DeleteVmsSize
// ------------------------------------------------------------------------------
func (a *HttpAppHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {

	var returnData ResponseMsg = ResponseMsg{}
	var result int = SUCCESS

	am.Applog.Print(2, "DelTest")
	returnData.Data = "DelTest"

	returnData.Result = result

	json.NewEncoder(w).Encode(returnData)
}
