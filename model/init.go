package model

import (
	"github.com/astaxie/beego/logs"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	//
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 分页参数
type PaginationParam struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// DB 数据库链接单例
var DB *gorm.DB

// Database 在中间件中初始化mysql链接
func Database(connString string) {
	db, err := gorm.Open("mysql", connString)
	if os.Getenv("GIN_MODE") == "debug" {
		db.LogMode(true)
	}else{
		db.LogMode(false)
	}
	// Error
	if err != nil {
		// util.Log().Panic("连接数据库不成功", err)
		logs.Error("连接数据库不成功", err)
	}
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(50)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	// 全局禁用表名复数
	db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//您可以通过定义DefaultTableNameHandler对默认表名应用任何规则。
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return os.Getenv("TABLE_PREFIX") + defaultTableName
	}

	DB = db
}
