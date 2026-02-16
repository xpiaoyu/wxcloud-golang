package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"

	"github.com/google/uuid"
)

func InsertRecordingHandler(w http.ResponseWriter, r *http.Request) {
	resp := &JsonResult{}

	// OpenId来自请求头，在微信云托管中，微信会自动将用户的openId放在请求头中，键名为"x-wx-openid"
	// 并且云托管中的openid是可信的，可以直接使用，不需要担心伪造问题
	openId := r.Header.Get("x-wx-openid")
	if strings.TrimSpace(openId) == "" {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 openId"
		outputJson(w, resp)
		return
	}

	fileId, err := getFileId(r)
	if err != nil {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 fileId"
		outputJson(w, resp)
		return
	}

	insertRecording(openId, fileId)
	resp.Code = 0
	outputJson(w, resp)
}

func getFileId(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	fileId, ok := body["fileId"]
	if !ok {
		return "", fmt.Errorf("缺少 fileId 参数")
	}

	return fileId.(string), nil
}

func insertRecording(openId string, fileId string) error {
	// 使用UUID作为录音记录的唯一标识符，确保每条记录都有一个独特的ID
	id := uuid.New().String()
	return dao.RecordingImp.InsertRecording(&model.RecordingModel{
		Id:        id,
		OpenId:    openId,
		FileId:    fileId,
		CreatedAt: time.Now().UTC(),
	})
}

func outputJson(w http.ResponseWriter, data *JsonResult) {
	msg, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("JSON序列化失败:", err)
		fmt.Fprintf(w, "内部错误")
		return
	}

	// 如果Code不为0，说明发生了错误，返回400状态码
	if data.Code != 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}
