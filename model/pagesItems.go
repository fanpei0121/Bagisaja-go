package model

type PagesItems struct {
	Id         int    `gorm:"primary_key" json:"id"`
	WidgetCode string `json:"widget_code"`
	PageCode   string `json:"page_code"`
	PositionId int    `json:"position_id"`
	Sort       int    `json:"sort"`
	Params     string `json:"params"`
}
