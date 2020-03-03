package model

import (
	"Bagisaja/util"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	OrderId     string  `gorm:"primary_key" json:"order_id"`
	GoodsAmount float64 `json:"goods_amount,omitempty"` // 商品总价
	OrderAmount float64 `json:"order_amount,omitempty"`
	PayStatus   int     `json:"pay_status,omitempty"`
	ShipStatus  int     `json:"ship_status,omitempty"`
	Status      int     `json:"status,omitempty"`
	OrderType   int     `json:"order_type,omitempty"`
	UserId      int     `json:"user_id,omitempty"`
	ShipAreaId  int     `json:"ship_area_id,omitempty"`
	ShipAddress string  `json:"ship_address,omitempty"`
	ShipName    string  `json:"ship_name,omitempty"`
	ShipMobile  string  `json:"ship_mobile,omitempty"`
	OrderPmt    float64 `json:"order_pmt"`
	GoodsPmt    float64 `json:"goods_pmt"`
	CouponPmt   float64 `json:"coupon_pmt"`
	Coupon      string  `json:"coupon"`
	Memo        string  `json:"memo"`
	Ip          string  `json:"ip,omitempty"`
	Ctime       int64   `json:"ctime,omitempty"`
	Isdel       int     `json:"isdel,omitempty"`
}

type OrderGoodsParam struct {
	ProductId int `json:"product_id"`
	Nums      int `json:"nums"`
}

type PmtParam struct {
	Goods      []OrderGoodsParam `json:"goods"`
	CouponCode string            `json:"coupon_code"` // 优惠券券码
}

// 下订单传的参数
type OrderParam struct {
	Goods       []OrderGoodsParam `json:"goods"`
	ShipAreaId  int               `json:"ship_area_id"`
	ShipAddress string            `json:"ship_address"`
	ShipName    string            `json:"ship_name"`
	ShipMobile  string            `json:"ship_mobile"`
	Memo        string            `json:"memo"`
	CouponCode  string            `json:"coupon_code"` // 优惠券券码
}

// 创建订单
func (t *Order) CreateOrder(param OrderParam, userId int, ip string) (string, error) {
	var addParam Order
	addParam.OrderId = util.BuildOrderNo() // 订单号
	addParam.PayStatus = 1                 // 未付款
	addParam.ShipStatus = 1                // 未发货
	addParam.Status = 1                    // 正常
	addParam.OrderType = 1                 // 普通订单
	addParam.UserId = userId
	addParam.ShipAreaId = param.ShipAreaId
	addParam.ShipAddress = param.ShipAddress
	addParam.ShipMobile = param.ShipMobile
	addParam.ShipName = param.ShipName
	addParam.Memo = param.Memo
	addParam.Ctime = time.Now().Unix()
	addParam.Ip = ip
	addParam.Coupon = param.CouponCode

	// 请注意，事务一旦开始，你就应该使用 tx 作为数据库句柄
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return "", err
	}
	user, _ := new(User).Info(userId)
	goodsPmtTotal := 0.000
	orderPmtTotal := 0.000
	couponPmtTotal := 0.000 // 优惠券优惠金额
	for _, vv := range param.Goods {
		productInfo := new(Products).GetInfoById(vv.ProductId)
		addParam.GoodsAmount += productInfo.Price * float64(vv.Nums) // 算出商品总价，用于计算订单打折优惠
	}
	for _, v := range param.Goods {
		goodsPmt, orderPmt, couponPmt, err := new(OrderItems).CreateItems(v, addParam, tx, *user, param.CouponCode)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		goodsPmtTotal += goodsPmt
		orderPmtTotal += orderPmt
		couponPmtTotal += couponPmt
	}
	addParam.GoodsPmt = goodsPmtTotal                                                            // 商品优惠金额
	addParam.OrderPmt = orderPmtTotal                                                            // 订单总优惠金额
	addParam.CouponPmt = couponPmtTotal                                                          // 优惠券优惠总金额
	addParam.OrderAmount = addParam.GoodsAmount - goodsPmtTotal - orderPmtTotal - couponPmtTotal // 实际销售金额
	if err := tx.Create(&addParam).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	return addParam.OrderId, tx.Commit().Error
}

// 根据购物车获取促销优惠金额
func (t *Order) GetPmtByGoods(pmtParam PmtParam, userId int) (float64, float64, float64) {
	var addParam Order
	addParam.OrderId = util.BuildOrderNo() // 订单号
	addParam.PayStatus = 1                 // 未付款
	addParam.ShipStatus = 1                // 未发货
	addParam.Status = 1                    // 正常
	addParam.OrderType = 1                 // 普通订单
	addParam.UserId = userId
	addParam.Ctime = time.Now().Unix()
	user, _ := new(User).Info(userId)
	goodsPmtTotal := 0.000
	orderPmtTotal := 0.000
	couponPmtTotal := 0.000
	for _, vv := range pmtParam.Goods {
		productInfo := new(Products).GetInfoById(vv.ProductId)
		addParam.GoodsAmount += productInfo.Price * float64(vv.Nums) // 算出商品总价，用于计算订单打折优惠
	}
	for _, v := range pmtParam.Goods {
		goodsPmt, orderPmt, couponPmt := new(OrderItems).GetPmtByOrderItem(v, addParam, *user, pmtParam.CouponCode)
		goodsPmtTotal += goodsPmt
		orderPmtTotal += orderPmt
		couponPmtTotal += couponPmt
	}
	return goodsPmtTotal, orderPmtTotal, couponPmtTotal
}

// 根据购物车获取用户能用的优惠券
func (t *Order) GetCouponByCart(cart []OrderGoodsParam, userId int) []CouponResponse {
	var addParam Order
	addParam.OrderId = util.BuildOrderNo() // 订单号
	addParam.PayStatus = 1                 // 未付款
	addParam.ShipStatus = 1                // 未发货
	addParam.Status = 1                    // 正常
	addParam.OrderType = 1                 // 普通订单
	addParam.UserId = userId
	addParam.Ctime = time.Now().Unix()
	user, _ := new(User).Info(userId)
	for _, vv := range cart {
		productInfo := new(Products).GetInfoById(vv.ProductId)
		addParam.GoodsAmount += productInfo.Price * float64(vv.Nums) // 算出商品总价，用于计算订单打折优惠
	}
	var promotionModel Promotion
	var productModel Products
	var goodsModel Goods
	var coupons []CouponResponse
	for _, v := range cart {
		product := productModel.GetInfoById(v.ProductId)
		goods := goodsModel.GetInfoById(product.GoodsId)
		var orderItem OrderItems
		orderItem.Nums = v.Nums
		res := promotionModel.GetCoupon(*goods, *user, orderItem, addParam)
		coupons = append(coupons, res...)
	}
	return coupons
}

type OrderResponse struct {
	Order
	CtimeFormat string               `json:"ctime_format"`
	OrderItems  []OrderItemsResponse `json:"order_items"`
}

// 订单详情
func (t *Order) GetInfoById(orderId string) *OrderResponse {
	var order OrderResponse
	if err := DB.Table("jshop_order").Where("order_id = ? AND isdel IS NULL", orderId).First(&order).Error; err != nil {
		return nil
	}
	order.CtimeFormat = util.TimeFormat(order.Ctime)
	order.OrderItems = new(OrderItems).GetItemsByOrderId(orderId)
	return &order
}

type OrderListParam struct {
	PaginationParam
	OrderType int `json:"order_type"`
}

// 订单列表
// orderType 0=全部 1待付款 2待发货 3待收货 4待评价
func (t *Order) GetListByOrderType(listParam OrderListParam, userId int) []OrderResponse {
	var orders []OrderResponse
	var db *gorm.DB
	switch listParam.OrderType {
	case 1:
		db = DB.Where("pay_status =  1")
	case 2:
		db = DB.Where("ship_status =  1")
	case 3:
		db = DB.Where("ship_status = 3")
	case 4:
		db = DB.Where("is_comment = 1")
	default:
		db = DB
	}
	offset := (listParam.Page - 1) * listParam.PageSize
	if err := db.Table("jshop_order").Select("order_id,order_amount,status,pay_status,ship_status,ctime").Where("user_id = ? AND isdel IS NULL", userId).Offset(offset).Limit(listParam.PageSize).Order("ctime desc").Scan(&orders).Error; err != nil {
		return nil
	}
	var orderItems OrderItems
	for k, order := range orders {
		orders[k].CtimeFormat = util.TimeFormat(order.Ctime)
		orders[k].OrderItems = orderItems.GetItemsByOrderId(order.OrderId)
	}
	return orders
}

// 删除订单
func (t *Order) Delete(orderId string, userId int) error {
	if err := DB.Table("jshop_order").Where("order_id = ? AND user_id = ? AND status <> 1", orderId, userId).Update("isdel", 1).Error; err != nil {
		logs.Error(err)
		return err
	}
	return nil
}
