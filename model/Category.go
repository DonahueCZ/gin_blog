package model

import (
	"errors"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
	"log"
)

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

//查询分类是否存在

func CheckCategory(name string) (code int) {
	var cate Category
	//这段代码，以id来排序，where来设定查询的限制，first寻找第一个匹配的记录
	db.Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED //2001分类存在
	}
	return errmsg.SUCCESS
}

// 新增分类

func CreateCate(data *Category) int {
	err := db.Create(&data).Error
	if err != nil {
		log.Printf("Error CreateUser: %v", err)
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//todo 查询分类下的所有文章

//查询分类列表

func GetCate(pageSize int, pageNum int) ([]Category, int64) {
	var cate []Category
	var total int64
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cate).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，可以返回空的分类列表而不是nil
			return []Category{}, 0
		}
		// 如果有其他错误，记录日志并处理错误
		log.Printf("Error GetCate: %v", err)
		return nil, 0
	}
	return cate, total
}

//编辑用户(密码以外的信息，密码单独做一个功能)

func EditCate(id int, data *Category) int {
	var cate Category
	// 首先检查用户ID是否存在
	err := db.First(&cate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果用户不存在，返回错误码
			return errmsg.ERROR_CATENAME_USED
		}
		// 如果查询过程中出现其他错误，也返回错误码
		return errmsg.ERROR
	}
	// 如果用户存在，继续更新操作
	var maps = make(map[string]interface{})
	maps["name"] = data.Name
	err = db.Model(&cate).Where("id= ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除分类

func DeleteCate(id int) int {
	var cate Category
	err = db.Where("id = ?", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
