package httpapp

import (
	"encoding/json"
	"net/http"
)

// ------------------------------------------------------------------------------
// DeleteVmsSize
// ------------------------------------------------------------------------------
func (a *HttpAppHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {

	var returnData ResponseMsg = ResponseMsg{}
	var result int = SUCCESS

	returnData.Result = result

	json.NewEncoder(w).Encode(returnData)
}
