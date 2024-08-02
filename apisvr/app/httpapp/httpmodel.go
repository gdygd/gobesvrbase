package httpapp

const (
	SUCCESS = 0
	FAIL    = 1

	BE_TIMEOUT = 4
)

// 쿠키만료시간
const LOGIN_COOKI_EMP_TM = 60 * 24 // 24 hour

const sysenvini = "./sys_env.ini"

// ---------------------------------------------------------------------------
// SessionObj
// ---------------------------------------------------------------------------
type SessionObj struct {
	Key int
}

// ---------------------------------------------------------------------------
// Response Msg
// ---------------------------------------------------------------------------
type ResponseMsg struct {
	Result  int         `json:"result"`
	Data    interface{} `json:"data"`
	ReqData interface{} `json:"reqdata"`
}
