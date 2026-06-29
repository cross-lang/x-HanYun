package constants

type BaseResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Debug string `json:"debug,omitempty"`
}
