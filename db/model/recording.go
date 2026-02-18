package model

import "time"

type RecordingModel struct {
	Id        string    `gorm:"column:id" json:"id"`
	OpenId    string    `gorm:"column:openId" json:"openId"`
	FileId    string    `gorm:"column:fileId" json:"fileId"`
	Duration  int       `gorm:"column:duration" json:"duration"` // 录音时长，单位为毫秒
	FileSize  int       `gorm:"column:fileSize" json:"fileSize"` // 录音文件大小，单位为字节
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	Timestamp int64     `gorm:"column:timestamp" json:"timestamp"` // 辅助排序字段
}
