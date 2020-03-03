package model

import (
	"github.com/astaxie/beego/logs"
	"strings"
	"time"
)

type GoodsComment struct {
	Id      int    `gorm:"primary_key" json:"id"`
	Score   int    `json:"score"`
	UserId  int    `json:"user_id"`
	GoodsId int    `json:"goods_id"`
	OrderId int    `json:"order_id"`
	Images  string `json:"images"`
	Content string `json:"content"`
	Display int    `json:"display"`
	Ctime   int64  `json:"ctime"`
}

// 用户评论
func (t *GoodsComment) Add(param GoodsComment) error {
	param.Ctime = time.Now().Unix()
	if param.Score < 1 {
		param.Score = 1
	}
	if param.Score > 5 {
		param.Score = 5
	}
	err := DB.Create(&param).Error
	if err != nil {
		logs.Error(err)
	}
	return err
}

type GoodsCommentResponse struct {
	GoodsComment
	Nickname   string   `json:"nickname"`
	Avatar     string   `json:"avatar"`
	ImageSlice []string `json:"image_slice"`
}

// 商品评论列表
func (t *GoodsComment) List(goodsId int, page int, pageSize int) []GoodsCommentResponse {
	var list []GoodsCommentResponse
	offet := (page - 1) * pageSize
	re := DB.Table("jshop_goods_comment a").
		Select("a.*,b.nickname,b.avatar").
		Joins("inner join jshop_user b on a.user_id = b.id").
		Where("a.goods_id = ? AND a.display = 1", goodsId).
		Offset(offet).
		Limit(pageSize).
		Find(&list).RecordNotFound()
	if re {
		return nil
	}
	for k, v := range list {
		list[k].ImageSlice = strings.Split(v.Images, ",")
	}
	return list
}
