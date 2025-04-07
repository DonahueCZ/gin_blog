package model

import (
	"fmt"
	"ginblog/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB
var err error

// 连接数据库,里面包含了一些基础的数据

func InitDb() {
	//后面是gorm的一个固定格式
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	)
	//初始化gorm框架与数据库连接， &gorm.Config{}是配置对象，初始化为默认配置
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能，它是 GORM 库底层使用的数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}

	//AutoMigrate 会创建表、缺失的外键、约束、列和索引。 如果大小、精度、是否为空可以更改，则 AutoMigrate 会改变列的类型。 出于保护您数据的目的，它 不会 删除未使用的列
	db.AutoMigrate(&User{}, &Article{}, &Category{})

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	//用来关闭数据库
	//sqlDB.Close()
}
