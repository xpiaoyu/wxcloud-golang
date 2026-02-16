package model

import "time"

type RecordingModel struct {
	Id        string    `gorm:"column:id" json:"id"`
	OpenId    string    `gorm:"column:openId" json:"openId"`
	FileId    string    `gorm:"column:fileId" json:"fileId"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	Timestamp int64     `gorm:"column:timestamp" json:"timestamp"` // 辅助排序字段
}
