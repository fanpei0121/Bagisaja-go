package api

import (
	"Bagisaja/model"
	"github.com/gin-gonic/gin"
)

type CouponController struct {
}

// 优惠券中心
func (t *CouponController) CouponCenter(c *gin.Context) {
	coupons := new(model.Coupon).CouponCenter()
	c.JSON(200, coupons)
}

// 用户领取优惠券
func (t *CouponController) ReceiveCoupon(c *gin.Context) {
	var coupon model.Coupon
	if err := c.ShouldBind(&coupon); err != nil || coupon.PromotionId <= 0 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	if err := new(model.Coupon).UserReceiveCoupon(userId, coupon.PromotionId); err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}
