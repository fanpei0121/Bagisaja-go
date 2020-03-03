package api

import (
	"Bagisaja/model"
	"Bagisaja/util"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"html"
	"os"
	"strconv"
)

type UserController struct {
}

// 用户短信验证码登录或注册
func (t *UserController) Login(c *gin.Context) {
	param := make(map[string]string)
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse("ParamError"))
		return
	}
	mobile := param["mobile"]
	mobile, ok := util.CheckIndonesiaMobile(mobile)
	if !ok {
		c.JSON(200, ErrorResponse("mobileError"))
		return
	}
	if os.Getenv("GIN_MODE") == "debug" { //测试环境
		if param["code"] != "111111" {
			c.JSON(200, ErrorResponse("codeError"))
			return
		}
	} else {
		// 验证redis的code码 暂时没做
	}
	token, err := new(model.User).RegisterOrLogin(mobile)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	ret := make(map[string]string)
	ret["token"] = token
	c.JSON(200, SuccessResponse("success", ret))
}

// 获取用户信息
func (t *UserController) Info(c *gin.Context) {
	uid, err := GetUidByHead(c)
	info, err := new(model.User).Info(uid)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	c.JSON(200, SuccessResponse("success", info))
}

var userShipModel model.UserShip

// 添加用户地址
func (t *UserController) AddUserShip(c *gin.Context) {
	var param model.UserShip
	if err := c.ShouldBind(&param); err != nil {
		logs.Error(err)
		c.JSON(200, ErrorResponse("ParamError"))
		return
	}
	if param.Name == "" || param.Address == "" || param.Mobile == "" || param.AreaId <= 0 || param.IsDef < 1 {
		c.JSON(200, ErrorResponse("ParamError"))
		return
	}
	param.UserId, _ = GetUidByHead(c)
	if err := userShipModel.Add(param); err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 用户地址信息
func (t *UserController) UserShipInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	info := userShipModel.Info(id)
	c.JSON(200, SuccessResponse(Success, info))
}

// 修改用户地址
func (t *UserController) UserShipUpdate(c *gin.Context) {
	var param model.UserShip
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if param.Name == "" || param.Address == "" || param.Mobile == "" || param.AreaId <= 0 || param.IsDef < 1 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	param.UserId, _ = GetUidByHead(c)
	if err := userShipModel.Update(param); err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 用户地址删除
func (t *UserController) UserShipDelete(c *gin.Context) {
	var param model.UserShip
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	if err := userShipModel.Delete(param.Id); err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 用户地址列表
func (t *UserController) UserShipList(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	list := userShipModel.List(userId)
	c.JSON(200, SuccessResponse(Success, list))
}

var CartModel model.Cart

// 货品加入购物车
func (t *UserController) CartAdd(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	var param model.CartParam
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if param.ProductId <= 0 || param.Nums <= 0 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if err := CartModel.Add(userId, param); err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 购物车列表
func (t *UserController) Carts(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	list := CartModel.List(userId)
	c.JSON(200, SuccessResponse(Success, list))
}

// 购物车删除
func (t *UserController) CartDelete(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	param := make(map[string]int)
	if err := c.ShouldBind(&param); err != nil || param["product_id"] <= 0 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if err := CartModel.Delete(userId, param["product_id"]); err != nil {
		logs.Error(err)
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 用户评价商品
func (t *UserController) GoodsCommentAdd(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	var param model.GoodsComment
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if param.GoodsId <= 0 || param.OrderId <= 0 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if param.Content == "" && param.Images == "" {
		c.JSON(200, ErrorResponse("ParamEmpty"))
		return
	}
	param.UserId = userId
	param.Content = html.EscapeString(param.Content) // 转义特殊字符
	param.Images = html.EscapeString(param.Images)
	err := new(model.GoodsComment).Add(param)
	if err != nil {
		c.JSON(200, ErrorResponse(Error))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}

// 会员签到
func (t *UserController) SignIn(c *gin.Context) {
	userId, _ := GetUidByHead(c)
	if err := new(model.UserPointLog).SignIn(userId); err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	c.JSON(200, SuccessResponse(Success))
}
