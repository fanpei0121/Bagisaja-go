package api

import (
	"Bagisaja/model"
	"github.com/gin-gonic/gin"
)

type AreaController struct {
}

var areaModel model.AreaIndonesia

// 获取所有省
func (t *AreaController) Provinces(c *gin.Context) {
	provinces, err := areaModel.GetAllProvince()
	if err != nil {
		c.JSON(200, ErrorResponse("error"))
	}
	c.JSON(200, SuccessResponse("success", provinces))
}

// 获取市
func (t *AreaController) Citys(c *gin.Context) {
	province := c.Query("province")
	if province == "" {
		c.JSON(200, ErrorResponse("ParamError"))
		return
	}
	citys := areaModel.GetCity(province)
	c.JSON(200, SuccessResponse("success", citys))
}

// 获取区
func (t *AreaController) Areas(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(200, ErrorResponse("ParamError"))
		return
	}
	areas := areaModel.GetArea(city)
	c.JSON(200, SuccessResponse("success", areas))
}
