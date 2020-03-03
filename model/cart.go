package model

import "github.com/astaxie/beego/logs"

type Cart struct {
	Id        int `gorm:"primary_key" json:"id"`
	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
	Nums      int `json:"nums"`
}

type CartParam struct {
	ProductId int `json:"product_id"`
	Nums      int `json:"nums"`
}

// 加入购物车
func (t *Cart) Add(userId int, param CartParam) error {
	cart := Cart{
		UserId:    userId,
		ProductId: param.ProductId,
		Nums:      param.Nums,
	}
	err := DB.Create(&cart).Error
	if err != nil {
		logs.Error(err)
	}
	return err
}

type ListResponse struct {
	ProductsResponse
	Nums      int `json:"nums"`
	ProductId int `json:"product_id"`
}

// 列表
func (t *Cart) List(userId int) []ListResponse {
	var list []ListResponse
	recordNotFound := DB.Table("jshop_cart a").
		Select("a.nums,a.product_id,b.*,c.url").
		Joins("inner join jshop_products b on a.product_id = b.id").
		Joins("left join jshop_images c on b.image_id = c.id").
		Where("a.user_id = ?", userId).
		Order("a.id desc").
		Find(&list).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return list
}

// 购物车删除
func (t *Cart) Delete(userId int, productId int) error {
	err := DB.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&Cart{}).Error
	return err
}
