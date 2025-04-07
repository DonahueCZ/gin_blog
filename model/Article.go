package model

import (
	"errors"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:ID"` //外键
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Cid     uint   `gorm:"type:int;not null" json:"cid"`
	Desc    string `gorm:"type:varchar(255)" json:"desc"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
}

// 新增文章

func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//查询分类下的所有文章

func GetCateArticle(id int, pageSize int, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("cid=?", id).Find(&cateArtList).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return cateArtList, errmsg.SUCCESS, total
}

// 查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err := db.Preload("Category").Where("id = ?", id).First(&art).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCESS
}

//查询文章列表

func GetArt(pageSize int, pageNum int) ([]Article, int, int64) {
	var articleList []Article
	var total int64
	//Preload预加载
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&articleList).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCESS, total
}

//编辑用户(密码以外的信息，密码单独做一个功能)

func EditArt(id int, data *Article) int {
	var art Article
	// 首先检查用户ID是否存在
	err := db.First(&art, id).Error
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
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img
	err = db.Model(&art).Where("id= ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除分类

func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ?", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
