package router

import (
	"Bagisaja/api"
	"Bagisaja/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Static("/static", "./static")
	// 中间件, 顺序不能改
	router.Use(middleware.Cors())
	router.POST("/upload", api.Upload) //  上传图片

	userController := new(api.UserController)
	orderController := new(api.OrderController)
	couponController := new(api.CouponController)
	router.POST("/user/login", userController.Login) // 登录
	user := router.Group("/user").Use(middleware.CheckLogin())
	{
		user.GET("info", userController.Info)                  // 账号信息
		user.POST("shipAdd", userController.AddUserShip)       // 添加用户地址
		user.GET("shipInfo", userController.UserShipInfo)      // 地址详情
		user.POST("shipUpdate", userController.UserShipUpdate) // 地址修改
		user.POST("shipDelete", userController.UserShipDelete) // 地址删除
		user.GET("ships", userController.UserShipList)         // 地址列表

		user.POST("cartAdd", userController.CartAdd)       // 货品加入购物车
		user.GET("carts", userController.Carts)            // 购物车列表
		user.POST("cartDelete", userController.CartDelete) // 购物车删除

		user.POST("commentAdd", userController.GoodsCommentAdd) // 用户评论

		user.POST("orderCreate", orderController.CreateOrder) // 用户下单
		user.POST("getPmt", orderController.GetPmt)           // 获取购物车的优惠信息
		user.POST("getCoupons", orderController.GetCoupons)   // 根据购物车获取用户能用的优惠券
		user.GET("orderInfo", orderController.OrderInfo)      // 订单详情
		user.POST("orderList", orderController.OrderList)     // 订单列表
		user.POST("orderDel", orderController.OrderDel)       // 订单删除

		user.POST("signIn", userController.SignIn) // 用户签到

		user.GET("couponCenter", couponController.CouponCenter)    // 待领取优惠券列表
		user.POST("receiveCoupon", couponController.ReceiveCoupon) // 用户领取优惠券
	}

	areaController := new(api.AreaController)
	router.GET("/provinces", areaController.Provinces) // 获取省
	router.GET("/citys", areaController.Citys)         // 获取市
	router.GET("/areas", areaController.Areas)         // 获取区

	goodsController := new(api.GoodsController)
	router.GET("/cats", goodsController.Cats)           // 获取分类
	router.GET("/goods", goodsController.GoodsList)     // 获取商品列表
	router.GET("/goodsInfo", goodsController.GoodsInfo) // 获取商品详情
	router.GET("/product", goodsController.Product)     // 获取规格详情
	router.POST("/comments", goodsController.Comments)  // 商品评论

	homeController := new(api.HomeController)
	router.GET("/home", homeController.Home) // 获取首页布局

	return router
}
