package model

import (
	"github.com/astaxie/beego/logs"
	"strconv"
)

type Setting struct {
	Skey  string `json:"skey"`
	Value string `json:"value"`
}

type SettingResponse struct {
	SignPointType            int `json:"sign_point_type"`            // 签到奖励类型 1=固定奖励,2=随机奖励
	FirstSignPoint           int `json:"first_sign_point"`           // 首次奖励积分
	ContinuitySignAdditional int `json:"continuity_sign_additional"` // 连续签到追加
	SignMostPoint            int `json:"sign_most_point"`            // 单日最大奖励
	SignRandomMin            int `json:"sign_random_min"`            // 随机奖励积分最小值
	SignRandomMax            int `json:"sign_random_max"`            // 随机奖励积分最大值
}

// 获取设置
func (t *Setting) GetSetting() *SettingResponse {
	var settings []Setting
	if err := DB.Find(&settings).Error; err != nil {
		logs.Error(err)
		return nil
	}
	var settingResponse SettingResponse
	var err error
	for _, v := range settings {
		switch v.Skey {
		case "sign_point_type":
			settingResponse.SignPointType, err = strconv.Atoi(v.Value)
		case "first_sign_point":
			settingResponse.FirstSignPoint, err = strconv.Atoi(v.Value)
		case "continuity_sign_additional":
			settingResponse.ContinuitySignAdditional, err = strconv.Atoi(v.Value)
		case "sign_most_point":
			settingResponse.SignMostPoint, err = strconv.Atoi(v.Value)
		case "sign_random_min":
			settingResponse.SignRandomMin, err = strconv.Atoi(v.Value)
		case "sign_random_max":
			settingResponse.SignRandomMax, err = strconv.Atoi(v.Value)

		}
		if err != nil {
			logs.Error(err)
			return nil
		}
	}
	return &settingResponse
}
