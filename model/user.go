package model

import (
	"Bagisaja/service/jwtx"
	"errors"
	"github.com/astaxie/beego/logs"
)

type User struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Sex      int    `json:"sex"`
	Birthday string `json:"birthday"`
	Grade    int    `json:"grade"` // 用户等级
	Status   int    `json:"status"`
	Point    int    `json:"point"`
}

// 注册或者登录
func (t *User) RegisterOrLogin(mobile string) (string, error) {
	var user User
	var uid int
	if DB.Where("mobile = ? AND isdel IS NULL", mobile).First(&user).RecordNotFound() { // 注册
		newUser := User{
			Mobile: mobile,
		}
		if err := DB.Create(&newUser).Error; err != nil {
			return "", err
		}
		uid = newUser.Id
	}
	if user.Status != 1 { // 账号被停用
		return "", errors.New("statusError")
	}
	uid = user.Id
	token, err := t.GetToken(uid)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 根据uid获取token
func (t *User) GetToken(id int) (string, error) {
	jwtParam := make(map[string]interface{})
	jwtParam["uid"] = id
	tokenString, err := jwtx.GenToken(jwtParam)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	return tokenString, nil
}

// 获取用户信息
func (t *User) Info(uid int) (*User, error) {
	var user User
	if DB.First(&user, uid).RecordNotFound() {
		return nil, errors.New("NotFound")
	}
	return &user, nil
}
