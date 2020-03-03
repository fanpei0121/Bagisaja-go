package model

import (
	"github.com/astaxie/beego/logs"
	"time"
)

type UserShip struct {
	Id      int    `gorm:"primary_key" json:"id"`
	UserId  int    `json:"user_id"`
	AreaId  int    `json:"area_id"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Utime   int64  `json:"utime"`
	IsDef   int    `json:"is_def"`
}

// 添加用户地址
func (t *UserShip) Add(param UserShip) error {
	var count int

	if err := DB.Model(&UserShip{}).Where("user_id = ?", param.UserId).Count(&count).Error; err != nil {
		logs.Error(err)
		return err
	}
	if count < 1 {
		param.IsDef = 1
	} else {
		if param.IsDef == 1 { // 默认地址
			if err := DB.Model(&UserShip{}).Where("user_id = ?", param.UserId).UpdateColumn("is_def", 2).Error; err != nil {
				return err
			}
		} else {
			param.IsDef = 2
		}
	}
	param.Utime = time.Now().Unix()
	err := DB.Create(&param).Error
	return err
}

// 地址信息
func (t *UserShip) Info(id int) *UserShip {
	var userShip UserShip
	if DB.First(&userShip, id).RecordNotFound() {
		return nil
	}
	return &userShip
}

// 地址修改
func (t *UserShip) Update(param UserShip) error {
	if param.IsDef == 1 { // 默认地址
		if err := DB.Model(&UserShip{}).Where("user_id = ?", param.UserId).UpdateColumn("is_def", 2).Error; err != nil {
			return err
		}
	} else {
		param.IsDef = 2
	}
	param.Utime = time.Now().Unix()
	err := DB.Model(&UserShip{}).Updates(&param).Error
	return err
}

// 地址删除
func (t *UserShip) Delete(id int) error {
	err := DB.Where("id = ?", id).Delete(UserShip{}).Error
	return err
}

// 地址列表
func (t *UserShip) List(userId int) []UserShip {
	var userShips []UserShip
	if DB.Where("user_id = ?", userId).Order("is_def asc").Find(&userShips).RecordNotFound() {
		return nil
	}
	return userShips
}
