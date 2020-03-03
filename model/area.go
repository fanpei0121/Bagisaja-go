package model

import (
	"github.com/astaxie/beego/logs"
)

type AreaIndonesia struct {
	Id       int    `gorm:"primary_key" json:"id"`
	PostCode string `json:"post_code"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	Street   string `json:"street"`
	Group    string `json:"group"`
}

type Area struct {
	Id         int    `gorm:"primary_key" json:"id"`
	ParentId   int    `json:"parent_id"`
	Depth      int    `json:"depth"`
	Name       string `json:"name"`
	PostalCode string `json:"postal_code"`
	Group      string `json:"group"`
	Sort       int    `json:"sort"`
}

// 获取所有省份
func (t *AreaIndonesia) GetAllProvince() ([]string, error) {
	var provinces []string
	sql := "SELECT distinct `province` FROM jshop_area_indonesia"
	err := DB.Raw(sql).Pluck("province", &provinces).Error
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return provinces, nil
}

// 获取市
func (t *AreaIndonesia) GetCity(province string) []string {
	var citys []string
	sql := "SELECT distinct `city` FROM jshop_area_indonesia WHERE `province` = ?"
	recordNotFound := DB.Raw(sql, province).Pluck("city", &citys).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return citys
}

// 获取区
func (t *AreaIndonesia) GetArea(city string) []AreaIndonesia {
	var areas []AreaIndonesia
	if err := DB.Select("area,street,post_code").Where("city = ?", city).Find(&areas).Error; err != nil {
		logs.Error(err)
		return nil
	}
	return areas
}
