package model

type Products struct {
	Id          int     `gorm:"primary_key" json:"id"`
	GoodsId     int     `json:"goods_id"`
	Sn          string  `json:"sn"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	FreezeStock int     `json:"freeze_stock"`
	ImageId     string  `json:"image_id"`
	SpesDesc    string  `json:"spes_desc"`
}

type ProductsResponse struct {
	Products
	Url  string `json:"url"`
	Name string `json:"name"` // 商品名称
	Bn   string `json:"bn"`   // 商品编码
}

// 获取规格详情
func (t *Products) GetInfo(goodsId int, spesDesc string) *ProductsResponse {
	var product ProductsResponse
	recordNotFound := DB.Table("jshop_products a").
		Select("a.*,b.url").
		Joins("left join jshop_images b on a.image_id = b.id").
		Where("a.goods_id = ? AND a.spes_desc = ?", goodsId, spesDesc).
		Scan(&product).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return &product
}

// 根据主键获取规格信息
func (t *Products) GetInfoById(productId int) *ProductsResponse {
	var product ProductsResponse
	recordNotFound := DB.Table("jshop_products a").
		Select("a.*,b.url,c.name,c.bn").
		Joins("left join jshop_images b on a.image_id = b.id").
		Joins("left join jshop_goods c on a.goods_id = c.id").
		Where("a.id = ?", productId).
		Scan(&product).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return &product
}
