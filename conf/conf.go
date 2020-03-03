package conf

import (
	"os"
	"Bagisaja/cache"
	"Bagisaja/model"

	"github.com/joho/godotenv"
	"github.com/astaxie/beego/logs"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	if err := godotenv.Load(".env"); err != nil {
		panic(".env load error")
	}

	// 设置日志
	initLogs()

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}

func initLogs()  {
	if os.Getenv("GIN_MODE") == "debug" {
		logs.SetLogger("console")
	}
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
}
