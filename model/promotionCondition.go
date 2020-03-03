package model

import (
	"Bagisaja/util"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type PromotionCondition struct {
	Id          int    `gorm:"primary_key" json:"id"`
	PromotionId int    `json:"promotion_id"`
	Code        string `json:"code"`
	Params      string `json:"params"`
}

// 获取促销条件列表
func (t *PromotionCondition) List(promotionId int) []PromotionCondition {
	var promotionCondition []PromotionCondition
	if DB.Where("promotion_id = ?", promotionId).Find(&promotionCondition).RecordNotFound() {
		return nil
	}
	return promotionCondition
}

// 检查是否满足条件
func (t *PromotionCondition) Check(goods Goods, user User, orderItem OrderItems, order Order, conditionList []PromotionCondition) bool {
	for _, v := range conditionList {
		paramsMap := make(map[string]string)
		if err := json.Unmarshal([]byte(v.Params), &paramsMap); err != nil {
			logs.Error(err)
			return false
		}
		switch v.Code {
		case "GOODS_IDS": // 判断指定某些商品满足
			idStr := strconv.Itoa(goods.Id)
			nums, err := strconv.Atoi(paramsMap["nums"])
			if err != nil {
				logs.Error(err)
				return false
			}
			goodsIdSlice := strings.Split(paramsMap["goods_id"], ",")
			_, isIn := util.InSlice(goodsIdSlice, idStr)
			if isIn == false || orderItem.Nums < nums {
				return false
			}
		case "GOODS_CATS": // 判断指定商品分类满足
			catId, err := strconv.Atoi(paramsMap["cat_id"])
			if err != nil {
				logs.Error(err)
				return false
			}
			nums, err := strconv.Atoi(paramsMap["nums"])
			if err != nil {
				logs.Error(err)
				return false
			}
			if goods.GoodsCatId != catId || orderItem.Nums < nums {
				return false
			}
		case "GOODS_BRANDS": // 指定商品品牌满足条件
			brandIdStr := strconv.Itoa(goods.BrandId)
			nums, err := strconv.Atoi(paramsMap["nums"])
			if err != nil {
				logs.Error(err)
				return false
			}
			brandIdSlice := strings.Split(paramsMap["brand_id"], ",")
			_, isIn := util.InSlice(brandIdSlice, brandIdStr)
			if isIn == false || orderItem.Nums < nums {
				return false
			}
		case "USER_GRADE": // 用户符合指定等级
			if user.Id > 0 {
				grades := util.GetKeys(paramsMap)
				_, isIn := util.InSlice(grades, strconv.Itoa(user.Grade))
				if isIn == false {
					return false
				}
			} else {
				return false
			}
		case "ORDER_FULL": // 订单满XX金额满足条件
			money, err := strconv.Atoi(paramsMap["money"])
			if err != nil {
				logs.Error(err)
				return false
			}
			if order.GoodsAmount < float64(money) {
				return false
			}

		}
	}
	return true
}
