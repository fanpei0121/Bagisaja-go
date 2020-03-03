package middleware

import (
	"Bagisaja/api"
	"github.com/gin-gonic/gin"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := api.GetUidByHead(c)
		if err != nil || uid <= 0 {
			c.JSON(401, "noLogin")
			c.Abort()
		}
		c.Next()
	}
}
