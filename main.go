package main

import (
	"Bagisaja/conf"
	"Bagisaja/router"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	// 装载路由
	r := router.NewRouter()
	r.Run(":9090")
}
