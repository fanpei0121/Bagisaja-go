package api

import (
	"Bagisaja/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type GoodsController struct {
}

var goodsCatModel model.GoodsCat

// 获取分类
func (t *GoodsController) Cats(c *gin.Context) {
	parentId, err := strconv.Atoi(c.DefaultQuery("parent_id", "0"))
	if err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	result := goodsCatModel.GetCats(parentId)
	c.JSON(200, SuccessResponse(Success, result))
}

var goodsModel model.Goods

// 商品列表
func (t *GoodsController) GoodsList(c *gin.Context) {
	catId, err := strconv.Atoi(c.Query("cat_id"))
	if err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	list := goodsModel.List(catId)
	c.JSON(200, SuccessResponse(Success, list))
}

// 商品详情
func (t *GoodsController) GoodsInfo(c *gin.Context) {
	goodsId, err := strconv.Atoi(c.Query("goods_id"))
	if err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	userId, _ := GetUidByHead(c)
	info := goodsModel.Info(goodsId, userId)
	c.JSON(200, SuccessResponse(Success, info))
}

// 获取规格详情
func (t *GoodsController) Product(c *gin.Context) {
	goodsId, err := strconv.Atoi(c.Query("goods_id"))
	if err != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	spesDesc := c.Query("spes_desc")
	if spesDesc == "" {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	info := new(model.Products).GetInfo(goodsId, spesDesc)
	c.JSON(200, SuccessResponse(Success, info))
}

type CommentParam struct {
	GoodsId  int `json:"goods_id"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// 商品评价列表
func (t *GoodsController) Comments(c *gin.Context) {
	var commentParam CommentParam
	if c.ShouldBind(&commentParam) != nil {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if commentParam.GoodsId <= 0 {
		c.JSON(200, ErrorResponse(ParamError))
		return
	}
	if commentParam.Page <= 0 {
		commentParam.Page = 1
	}
	if commentParam.PageSize <= 0 {
		commentParam.PageSize = 10
	}
	list := new(model.GoodsComment).List(commentParam.GoodsId, commentParam.Page, commentParam.PageSize)
	c.JSON(200, SuccessResponse(Success, list))
}
