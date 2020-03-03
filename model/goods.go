package model

import "github.com/techleeone/gophp/serialize"

type Goods struct {
	Id         int     `gorm:"primary_key" json:"id"`
	Name       string  `json:"name"`
	Brief      string  `json:"brief"`
	Mktprice   float64 `json:"mktprice"`
	Costprice  float64 `json:"costprice"`
	ImageId    string  `json:"image_id,omitempty"`
	GoodsCatId int     `json:"goods_cat_id"`
	BrandId    int     `json:"brand_id"`
	Intro      string  `json:"intro"`
	SpesDesc   string  `json:"spes_desc"`
}

// 列表数据
type GoodsListResponse struct {
	Goods
	Url string `json:"image_url"`
}

// 列表
func (t *Goods) List(catId int) []GoodsListResponse {
	var goods []GoodsListResponse
	if DB.Table("jshop_goods a").
		Select("a.*,b.url").
		Joins("left join jshop_images b on a.image_id=b.id").
		Where("a.goods_cat_id = ? AND a.marketable = 1 AND a.isdel IS NULL", catId).Order("sort asc,id desc").Find(&goods).RecordNotFound() {
		return nil
	}
	return goods
}

// 详情数据
type GoodsInfoResponse struct {
	Goods
	GoodsBanner   []string    `json:"goods_banner"`
	SpesDescMap   interface{} `json:"spes_desc_map"`
	PromotionList []Promotion `json:"promotion_list"`
}

// 详情
func (t *Goods) Info(id int, userId int) *GoodsInfoResponse {
	var goodsInfo GoodsInfoResponse
	recordNotFound := DB.Table("jshop_goods").
		Select("*").
		Where("id = ?", id).
		First(&goodsInfo.Goods).RecordNotFound()
	if recordNotFound {
		return nil
	}
	// unserialize() in php
	goodsInfo.SpesDescMap, _ = serialize.UnMarshal([]byte(goodsInfo.Goods.SpesDesc)) // 商品规格
	goodsInfo.GoodsBanner = new(GoodsImages).GetBannerByGoodsId(id)                  // 首图
	user, _ := new(User).Info(userId)
	var orderItem OrderItems
	orderItem.Nums = 999 // 假设用户会买多个商品，把促销信息显示出来
	var order Order
	order.GoodsAmount = 99999.000 // 假设用户订单会超过10000,把促销信息显示出来
	var promotionModel Promotion
	goodsInfo.PromotionList = promotionModel.GetPromotionByGoods(goodsInfo.Goods, *user, orderItem, order) // 促销列表
	return &goodsInfo
}

// 获取最简单的goods信息
func (t *Goods) GetInfoById(id int) *Goods {
	var goods Goods
	if DB.Where("id = ?", id).First(&goods).RecordNotFound() {
		return nil
	}
	return &goods
}

