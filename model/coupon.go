package model

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"time"
)

type Coupon struct {
	CouponCode  string `gorm:"primary_key" json:"coupon_code,omitempty"`
	PromotionId int    `json:"promotion_id,omitempty"`
	IsUsed      int    `json:"is_used,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	Ctime       int64  `json:"ctime,omitempty"`
	Utime       int64  `json:"utime,omitempty"`
}

// 根据用户id获取优惠券ids
func (t *Coupon) GePromotionIdsByUserId(userId int) []int {
	var promotionIds []int
	if err := DB.Table("jshop_coupon").Where("user_id = ? AND is_used = 1", userId).Pluck("promotion_id", &promotionIds).Error; err != nil {
		return nil
	}
	return promotionIds
}

type CouponResponse struct {
	Coupon
	Name           string `json:"name"`            // 优惠券名称
	IntegralNumber int    `json:"integral_number"` // 所需积分
}

// 根据优惠券列表获取具体用户领取的优惠券
func (t *Coupon) GetListByPromotions(promotions []Promotion, userId int) []CouponResponse {
	var coupons []CouponResponse
	var promotionIds []int
	for _, v := range promotions {
		promotionIds = append(promotionIds, v.Id)
	}
	err := DB.Table("jshop_coupon a").
		Select("a.*,b.name").
		Joins("inner join jshop_promotion b on a.promotion_id = b.id").
		Where("a.promotion_id in (?) AND user_id = ? AND is_used = 1", promotionIds, userId).
		Scan(&coupons).Error
	if err != nil {
		return nil
	}
	return coupons
}

// 根据券码获取优惠券信息
func (t *Coupon) GetPromotionByCouponCode(couponCode string) []Promotion {
	var promotion []Promotion
	err := DB.Table("jshop_coupon a").
		Select("b.*").
		Joins("inner join jshop_promotion b on a.promotion_id =b.id").
		Where("a.coupon_code = ?", couponCode).
		Scan(&promotion).Error
	if err != nil {
		return nil
	}
	return promotion
}

// 待用户领取优惠券列表
func (t *Coupon) CouponCenter() []CouponResponse {
	nowTimeUnix := time.Now().Unix()
	var couponResponse []CouponResponse
	err := DB.Table("jshop_coupon a").
		Select("a.promotion_id,b.*").
		Joins("inner join jshop_promotion b on a.promotion_id = b.id").
		Where("b.status = 1 AND b.type = 2 AND b.auto_receive = 3 AND b.stime <= ? AND b.etime > ? AND b.isdel IS NULL", nowTimeUnix, nowTimeUnix).
		Where("a.user_id IS NULL").
		Group("a.promotion_id").
		Scan(&couponResponse).Error
	if err != nil {
		return nil
	}
	return couponResponse
}

// 用户领取优惠券
func (t *Coupon) UserReceiveCoupon(userId int, promotionId int) error {
	var num int
	if err := DB.Table("jshop_coupon").Where("user_id = ? AND promotion_id = ?", userId, promotionId).Count(&num).Error; err != nil {
		logs.Error(err)
		return err
	}
	if num > 0 {
		return errors.New("alreadyOwned")
	}
	var coupon Coupon
	if err := DB.Where("user_id IS NULL AND promotion_id = ?", promotionId).First(&coupon).Error; err != nil {
		return errors.New("noCoupons")
	}
	if err := DB.Model(&coupon).Update("user_id", userId).Error; err != nil {
		return err
	}
	return nil
}
