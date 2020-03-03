package model

type GoodsImages struct {
	GoodsId int    `json:"goods_id"`
	ImageId string `json:"image_id"`
	Sort    int    `json:"sort"`
}

// 获取首图
func (t *GoodsImages) GetBannerByGoodsId(goodsId int) []string {
	var banners []string
	recordNotFound := DB.Table("jshop_goods_images a").
		Select("b.url").
		Joins("left join jshop_images b on a.image_id = b.id").
		Where("a.goods_id = ?", goodsId).
		Order("a.sort asc").
		Pluck("b.url", &banners).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return banners
}
