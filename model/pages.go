package model

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type Pages struct {
	Id   int    `gorm:"primary_key" json:"id"`
	Code string `json:"code"`
}

type PagesItemsResponse struct {
	PagesItems
	ParamMap map[string]interface{}
}

// 获取页面布局数据
func (t *Pages) GetHomePage(id int) []PagesItemsResponse {
	var pagesItems []PagesItemsResponse
	recordNotFound := DB.Table("jshop_pages a").
		Select("b.*").
		Joins("left join jshop_pages_items b on a.code=b.page_code").
		Where("a.id = ?", id).
		Order("b.position_id asc").
		Scan(&pagesItems).RecordNotFound()
	if recordNotFound {
		return nil
	}
	for k, v := range pagesItems {
		if err := json.Unmarshal([]byte(v.Params), &pagesItems[k].ParamMap); err != nil {
			logs.Error(err)
		}

	}
	return pagesItems
}
