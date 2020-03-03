package model

import (
	"Bagisaja/util"
	"errors"
	"github.com/astaxie/beego/logs"
	"time"
)

type UserPointLog struct {
	Id      int    `gorm:"primary_key" json:"id"`
	UserId  int    `json:"user_id"`
	Type    int    `json:"type"`
	Num     int    `json:"num"`
	Balance int    `json:"balance"`
	Remarks string `json:"remarks"`
	Ctime   int64  `json:"ctime"`
}

// 用户签到
func (t *UserPointLog) SignIn(userId int) error {
	settings := new(Setting).GetSetting()
	num := 0
	lastLog := t.GetLastSignInLog(userId)
	var lastTime string
	if lastLog == nil { // 没有积分记录
		lastTime = ""
	} else {
		lastTime = time.Unix(lastLog.Ctime, 0).Format("2006-01-02")
	}
	if lastTime == time.Now().Format("2006-01-02") {
		return errors.New("AlreadySignedIn")
	}
	if settings.SignPointType == 1 { //1=固定奖励
		continuityTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // 连续的时间
		if continuityTime == lastTime {                                     // 连续签到
			num = settings.FirstSignPoint + settings.ContinuitySignAdditional
		} else { // 不连续，获取首次签到积分
			num = settings.FirstSignPoint
		}
		if num > settings.SignMostPoint { // 连续签到超出最大获取积分
			num = settings.SignMostPoint
		}
	}
	if settings.SignPointType == 2 { // 随机奖励
		num = util.RangeNum(settings.SignRandomMin, settings.SignRandomMax)
	}
	user, _ := new(User).Info(userId)
	// 请注意，事务一旦开始，你就应该使用 tx 作为数据库句柄
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	newPoint := user.Point+num
	if err := tx.Model(&user).Update("point", newPoint).Error; err != nil {
		tx.Rollback()
		logs.Error(err)
		return err
	}
	var userPointLog UserPointLog
	userPointLog.UserId = userId
	userPointLog.Type = 1
	userPointLog.Num = num
	userPointLog.Balance = newPoint
	userPointLog.Ctime = time.Now().Unix()
	userPointLog.Remarks = "签到: " + user.Mobile + " 积分奖励"
	if err := tx.Create(&userPointLog).Error; err != nil {
		logs.Error(err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 获取最近一次的签到记录
func (t *UserPointLog) GetLastSignInLog(userId int) *UserPointLog {
	var userPoint UserPointLog
	if err := DB.Where("type = 1 AND user_id = ?", userId).Order("ctime desc").First(&userPoint).Error; err != nil {
		return nil
	}
	return &userPoint
}
