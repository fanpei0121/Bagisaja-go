package model

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type OrderItems struct {
	Id              int     `gorm:"primary_key" json:"id"`
	OrderId         string  `json:"order_id"`
	GoodsId         int     `json:"goods_id"`
	ProductId       int     `json:"product_id"`
	Sn              string  `json:"sn"`
	Bn              string  `json:"bn"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	ImageUrl        string  `json:"image_url"`
	Nums            int     `json:"nums"`
	Amount          float64 `json:"amount"`
	PromotionAmount float64 `json:"promotion_amount"`
	Addon           string  `json:"addon"` // 货品明细序列号存储
	Utime           int64   `json:"utime"`
	PromotionList   string  `json:"promotion_list"`
}

// 新增订单商品明细
func (t *OrderItems) CreateItems(param OrderGoodsParam, order Order, tx *gorm.DB, user User, couponCode string) (float64, float64, float64, error) {
	productInfo := new(Products).GetInfoById(param.ProductId)
	if productInfo.Stock-param.Nums < productInfo.FreezeStock {
		return 0, 0, 0, errors.New("noStock")
	}
	var addParam OrderItems
	goodsMoney := 0.000  //该商品的商品优惠金额
	orderMoney := 0.000  // 该商品的订单优惠金额
	couponMoney := 0.000 // 优惠券优惠金额
	if productInfo != nil {
		addParam.OrderId = order.OrderId
		addParam.GoodsId = productInfo.GoodsId
		addParam.ProductId = param.ProductId
		addParam.Sn = productInfo.Sn
		addParam.Bn = productInfo.Bn
		addParam.Name = productInfo.Name
		addParam.Price = productInfo.Price
		addParam.ImageUrl = productInfo.Url
		addParam.Nums = param.Nums
		addParam.Amount = productInfo.Price * float64(param.Nums)                                 // 总价
		goodsMoney, orderMoney, couponMoney = t.GetPmtByOrderItem(param, order, user, couponCode) // 优惠信息
		addParam.PromotionAmount = goodsMoney
		addParam.Addon = productInfo.SpesDesc
		addParam.Utime = time.Now().Unix()
		addParam.PromotionList = "test"
		if err := tx.Create(&addParam).Error; err != nil {
			logs.Error(err)
			return 0.000, 0.000, 0.000, err
		}
		// 减少库存
		if err := tx.Table("jshop_products").Where("id = ?", param.ProductId).Update("stock", gorm.Expr("stock - ?", param.Nums)).Error; err != nil {
			logs.Error(err)
			return 0, 0, 0, err
		}
	}
	return goodsMoney, orderMoney, couponMoney, nil
}

// 获取单个商品的优惠金额 couponCode优惠券码
func (t *OrderItems) GetPmtByOrderItem(param OrderGoodsParam, order Order, user User, couponCode string) (float64, float64, float64) {
	productInfo := new(Products).GetInfoById(param.ProductId)
	var addParam OrderItems
	goodsMoney := 0.000  // 该商品优惠金额
	orderMoney := 0.000  // 该商品的订单优惠金额
	couponMoney := 0.000 // 优惠券优惠金额
	if productInfo != nil {
		addParam.OrderId = order.OrderId
		addParam.GoodsId = productInfo.GoodsId
		addParam.ProductId = param.ProductId
		addParam.Sn = productInfo.Sn
		addParam.Bn = productInfo.Bn
		addParam.Name = productInfo.Name
		addParam.Price = productInfo.Price
		addParam.ImageUrl = productInfo.Url
		addParam.Nums = param.Nums
		addParam.Amount = productInfo.Price * float64(param.Nums) // 总价
		var promotionModel Promotion
		goodsInfo := new(Goods).GetInfoById(addParam.GoodsId)
		promotions := promotionModel.GetPromotionByGoods(*goodsInfo, user, addParam, order)
		goodsMoney, orderMoney = promotionModel.ToResultByPromotions(promotions, addParam, order) // 商品优惠总金额
		// 算优惠券金额
		couponPromotion := new(Coupon).GetPromotionByCouponCode(couponCode)                                     // 优惠券
		couponPromotion = promotionModel.GetAllowPromotions(couponPromotion, *goodsInfo, user, addParam, order) // 筛选满足条件的优惠券
		couGoodsPmt, couOrderPmt := promotionModel.ToResultByPromotions(couponPromotion, addParam, order)
		couponMoney, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", couGoodsPmt+couOrderPmt), 64)
	}
	return goodsMoney, orderMoney, couponMoney // 商品优惠总金额
}

type OrderItemsResponse struct {
	OrderItems
	Name string `json:"name"`
}

// 根据订单id获取商品
func (t *OrderItems) GetItemsByOrderId(orderId string) []OrderItemsResponse {
	var list []OrderItemsResponse
	err := DB.Table("jshop_order_items a").
		Select("a.*,b.name").
		Joins("left join jshop_goods b on a.goods_id = b.id").
		Where("a.order_id = ?", orderId).
		Scan(&list).Error
	if err != nil {
		return nil
	}
	return list
}
