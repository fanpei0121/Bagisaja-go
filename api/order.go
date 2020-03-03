package api

import (
	"Bagisaja/model"
	"Bagisaja/util"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
}

var orderModel model.Order

// 下单
func (t *OrderController) CreateOrder(c *gin.Context) {
	var param model.OrderParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	for _, v := range param.Goods {
		if v.ProductId <= 0 || v.Nums <= 0 {
			c.JSON(200, ErrorResponse(ParamError))
			return
		}
	}
	if param.ShipAreaId <= 0 || param.ShipName == "" || param.ShipMobile == "" || param.ShipAddress == "" {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	ip := util.RemoteIp(c.Request)
	orderId, err := orderModel.CreateOrder(param, userId, ip)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	res := make(map[string]string)
	res["order_id"] = orderId
	c.JSON(200, SuccessResponse(Success, res))
}

// 获取购物车的促销优惠信息
func (t *OrderController) GetPmt(c *gin.Context) {
	var param model.PmtParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	goodsPmtTotal, orderPmtTotal, couponPmtTotal := orderModel.GetPmtByGoods(param, userId)
	result := make(map[string]float64)
	result["goods_pmt_total"] = goodsPmtTotal
	result["order_pmt_total"] = orderPmtTotal
	result["coupon_pmt_total"] = couponPmtTotal // 优惠券优惠金额
	c.JSON(200, SuccessResponse(Success, result))
}

// 根据购物车获取用户领取的优惠券
func (t *OrderController) GetCoupons(c *gin.Context) {
	var param []model.OrderGoodsParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	coupons := new(model.Order).GetCouponByCart(param, userId)
	c.JSON(200, SuccessResponse(Success, coupons))
}

// 订单详情
func (t *OrderController) OrderInfo(c *gin.Context) {
	orderId := c.Query("order_id")
	if orderId == "" {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	orderResponse := new(model.Order).GetInfoById(orderId)
	if orderResponse.UserId != userId {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success, orderResponse))
}

// 订单列表
// orderType 0=全部 1待付款 2待发货 3待收货 4待评价
func (t *OrderController) OrderList(c *gin.Context) {
	var listParam model.OrderListParam
	if err := c.ShouldBind(&listParam); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if listParam.Page <= 0 {
		listParam.Page = 1
	}
	if listParam.PageSize <= 0 {
		listParam.PageSize = 10
	}
	userId, _ := GetUidByHead(c)
	orderResponse := new(model.Order).GetListByOrderType(listParam, userId)
	c.JSON(200, SuccessResponse(Success, orderResponse))
}

// 删除订单
func (t *OrderController) OrderDel(c *gin.Context) {
	param := make(map[string]string)
	if err := c.ShouldBind(&param); err != nil || param["order_id"] == "" {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	err := new(model.Order).Delete(param["order_id"], userId)
	if err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}
