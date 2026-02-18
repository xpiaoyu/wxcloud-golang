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

	// 从请求体中获取 fileId 参数，fileId 是用户上传的录音文件在微信云存储中的唯一标识符
	fileId, err := getFileId(r)
	if err != nil {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 fileId"
		outputJson(w, resp)
		return
	}
	if strings.TrimSpace(fileId) == "" {
		resp.Code = -1
		resp.ErrorMsg = "fileId 不能为空"
		outputJson(w, resp)
		return
	}

	duration, err := getDuration(r)
	if err != nil {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 duration"
		outputJson(w, resp)
		return
	}

	fileSize, err := getFileSize(r)
	if err != nil {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 fileSize"
		outputJson(w, resp)
		return
	}

	insertRecording(openId, fileId, duration, fileSize)
	resp.Code = 0
	outputJson(w, resp)
}

func GetRecordingsByOpenIdHandler(w http.ResponseWriter, r *http.Request) {
	resp := &JsonResult{}
	openId := r.Header.Get("x-wx-openid")
	if strings.TrimSpace(openId) == "" {
		resp.Code = -1
		resp.ErrorMsg = "未获取到 openId"
		outputJson(w, resp)
		return
	}
	recordings, err := dao.RecordingImp.GetRecordingsByOpenId(openId, 0)
	if err != nil {
		resp.Code = -1
		resp.ErrorMsg = "查询录音记录失败"
		outputJson(w, resp)
		return
	}
	resp.Code = 0
	resp.Data = recordings
	outputJson(w, resp)
}

func getFileId(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}

	fileId, ok := body["fileId"]
	if !ok {
		return "", fmt.Errorf("缺少 fileId 参数")
	}

	return fileId.(string), nil
}

func getDuration(r *http.Request) (int, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return 0, err
	}

	duration, ok := body["duration"]
	if !ok {
		return 0, fmt.Errorf("缺少 duration 参数")
	}

	return int(duration.(float64)), nil
}

func getFileSize(r *http.Request) (int, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return 0, err
	}

	fileSize, ok := body["fileSize"]
	if !ok {
		return 0, fmt.Errorf("缺少 fileSize 参数")
	}
	return int(fileSize.(float64)), nil
}

func insertRecording(openId string, fileId string, duration int, fileSize int) error {
	// 使用UUID作为录音记录的唯一标识符，确保每条记录都有一个独特的ID
	id := uuid.New().String()
	return dao.RecordingImp.InsertRecording(&model.RecordingModel{
		Id:        id,
		OpenId:    openId,
		FileId:    fileId,
		Duration:  duration,
		FileSize:  fileSize,
		CreatedAt: time.Now().UTC(),
		Timestamp: time.Now().UnixNano(),
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
