package service

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"` // 0表示成功，非0表示失败
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}
