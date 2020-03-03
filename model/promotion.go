package model

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

type Promotion struct {
	Id          int    `gorm:"primary_key" json:"id"`
	Name        string `json:"name"`
	Status      int    `json:"status"`
	Type        int    `json:"type"`
	Exclusive   int    `json:"exclusive"`
	AutoReceive int    `json:"auto_receive"`
	Params      string `json:"params"`
	Stime       int    `json:"stime"`
	Etime       int    `json:"etime"`
	Isdel       int    `json:"isdel"`
}

// 根据商品获取促销列表
func (t *Promotion) GetPromotionByGoods(goods Goods, user User, orderItem OrderItems, order Order) []Promotion {
	promotions := t.GetPromotions(1)
	var promotionResponse []Promotion
	promotionResponse = t.GetAllowPromotions(promotions, goods, user, orderItem, order)  // 筛选满足条件的促销列表
	return promotionResponse
}

// 根据用户和商品获取优惠券列表
func (t *Promotion) GetCoupon(goods Goods, user User, orderItem OrderItems, order Order) []CouponResponse {
	var couponModel Coupon
	promotionIds := couponModel.GePromotionIdsByUserId(user.Id)
	promotions := t.GetPromotionsByIds(promotionIds)
	allowPromotions := t.GetAllowPromotions(promotions, goods, user, orderItem, order)
	result := couponModel.GetListByPromotions(allowPromotions, user.Id)
	return result
}

// 筛选满足要求的促销或优惠卷
func (t *Promotion) GetAllowPromotions(promotions []Promotion, goods Goods, user User, orderItem OrderItems, order Order) []Promotion {
	var promotionResponse []Promotion
	promotionConditionModel := new(PromotionCondition)
	for _, v := range promotions {
		conditionList := promotionConditionModel.List(v.Id) // 条件列表
		re := promotionConditionModel.Check(goods, user, orderItem, order, conditionList)
		if re { // 该促销满足条件
			promotionResponse = append(promotionResponse, v)
		}
		if v.Exclusive == 2 { //如果排他，就跳出循环，不执行下面的促销了
			break
		}
	}
	return promotionResponse
}

// 根据type获取促销列表 1促销，2优惠券，3团购，4秒杀
func (t *Promotion) GetPromotions(typeId int) []Promotion {
	var promotions []Promotion
	nowTime := time.Now().Unix()
	recordNotFound := DB.Where("status = 1 AND type = ? AND stime <= ? AND etime > ?", typeId, nowTime, nowTime).Order("sort asc").Find(&promotions).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return promotions
}

// 根据protionId获取优惠券列表
func (t *Promotion) GetPromotionsByIds(promotionIds []int) []Promotion {
	var promotions []Promotion
	if err := DB.Where("id in (?)", promotionIds).Find(&promotions).Error; err != nil {
		return nil
	}
	return promotions
}

// 根据购物车商品计算商品优惠金额
/*func (t *Promotion) ToResult(orderItems OrderItems, user User, order Order) (float64, float64) {
	goodsInfo := new(Goods).GetInfoById(orderItems.GoodsId)
	promotions := t.GetPromotionByGoods(*goodsInfo, user, orderItems, order) // 满足条件的促销列表
	return t.ToResultByPromotions(promotions, orderItems, order)
}*/

// 根据促销或优惠券列表获取优惠信息
func (t *Promotion) ToResultByPromotions(promotions []Promotion, orderItems OrderItems, order Order) (float64, float64) {
	// 计算结果
	goodsMoney := 0.000 // 商品优惠金额
	orderMoney := 0.000 // 订单优惠金额
	for _, v := range promotions {
		results := new(PromotionResult).GetListByPromotionId(v.Id) // 结果列表
		for _, result := range results {
			paramsMap := make(map[string]string)
			if err := json.Unmarshal([]byte(result.Params), &paramsMap); err != nil {
				logs.Error(err)
				return 0.000, 0.000
			}
			switch result.Code {
			case "GOODS_REDUCE": // 指定商品减固定金额
				deMoney, err := strconv.Atoi(paramsMap["money"])
				if err != nil {
					logs.Error(err)
				}
				goodsMoney += float64(deMoney)
			case "GOODS_DISCOUNT": // 指定商品打X折
				discount, err := strconv.Atoi(paramsMap["discount"])
				if err != nil {
					logs.Error(err)
				}
				goodsMoney += orderItems.Price * (1 - float64(discount)*0.1)
			case "GOODS_ONE_PRICE": // 指定商品一口价
				deMoney, err := strconv.Atoi(paramsMap["money"])
				if err != nil {
					logs.Error(err)
				}
				goodsMoney += orderItems.Amount - float64(deMoney)
			case "GOODS_HALF_PRICE": // 指定商品每第几件减指定金额
				num, err := strconv.Atoi(paramsMap["num"])
				if err != nil {
					logs.Error(err)
				}
				deMoney, err := strconv.Atoi(paramsMap["money"])
				if err != nil {
					logs.Error(err)
				}
				if int(orderItems.Nums/num) > 0 {
					goodsMoney += float64(deMoney * int(orderItems.Nums/num))
				}
			case "ORDER_REDUCE": // 订单减指定金额
				money, err := strconv.Atoi(paramsMap["money"])
				if err != nil {
					logs.Error(err)
				}
				orderMoney += float64(money)
			case "ORDER_DISCOUNT": // 订单打X折
				discount, err := strconv.Atoi(paramsMap["discount"])
				if err != nil {
					logs.Error(err)
				}
				orderMoney += order.GoodsAmount * (1 - float64(discount)*0.1)
			default:
				goodsMoney += 0.000
				orderMoney += 0.000
			}
		}
	}
	return goodsMoney, orderMoney
}


