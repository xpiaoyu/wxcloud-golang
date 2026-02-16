package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OpenIdResult 返回结构
type OpenIdResult struct {
	OpenId   string `json:"openid"`
	ErrorMsg string `json:"errorMsg,omitempty"`
}

func GetOpenIdHandler(w http.ResponseWriter, r *http.Request) {
	openId := r.Header.Get("x-wx-openid")
	if openId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "未获取到openId")
		return
	}
	res := &OpenIdResult{}
	res.OpenId = openId
	msg, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
