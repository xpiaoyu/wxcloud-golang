package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "Counters"

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

// RecordingInterfaceImp 录音记录接口实现

func (imp *RecordingInterfaceImp) InsertRecording(recording *model.RecordingModel) error {
	cli := db.Get()
	return cli.Table("Recordings").Create(recording).Error
}

func (imp *RecordingInterfaceImp) GetRecordingsByOpenId(openId string, lastTimestamp int64) ([]*model.RecordingModel, error) {
	var err error
	var recordings []*model.RecordingModel
	cli := db.Get()
	if lastTimestamp > 0 {
		err = cli.Table("Recordings").Where("openId = ? AND timestamp < ?", openId, lastTimestamp).Order("timestamp DESC").Find(&recordings).Error
	} else {
		err = cli.Table("Recordings").Where("openId = ?", openId).Order("timestamp DESC").Find(&recordings).Error
	}
	return recordings, err
}
