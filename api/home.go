package api

import (
	"Bagisaja/model"
	"github.com/gin-gonic/gin"
)

type HomeController struct {
}

// 首页布局
func (t *HomeController) Home(c *gin.Context) {
	list := new(model.Pages).GetHomePage(1)
	c.JSON(200, SuccessResponse(Success, list))
}
