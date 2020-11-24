package model

import (
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB gorm.DB 对象
var DB *gorm.DB

// ConnectDB 初始化模型
func ConnectDB() *gorm.DB {

    var err error

    config := mysql.New(mysql.Config{
        DSN: "root:123456@tcp(localhost:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
    })

    // 准备数据库连接池
    DB, err = gorm.Open(config, &gorm.Config{})

    logger.LogError(err)

    return DB
}