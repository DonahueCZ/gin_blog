package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加分类

func AddCategory(c *gin.Context) {
	//
	var data model.Category
	_ = c.ShouldBind(&data)
	code = model.CheckCategory(data.Name)
	if code == errmsg.SUCCESS {
		model.CreateCate(&data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		code = errmsg.ERROR_CATENAME_USED
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

//todo 查询分类下的所有文章

// 查询分类列表

func GetCate(c *gin.Context) {
	// 分页功能
	//queryh返回的是string类型，但是我们需要int，用strconv.Atoi转换
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	//
	if pageSize <= 0 {
		pageSize = 10 // 默认页面大小
	}
	if pageNum <= 0 {
		pageNum = 1 // 默认页码
	}
	data, total := model.GetCate(pageSize, pageNum)
	code = errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

//编辑分类名

func EditCate(c *gin.Context) {
	//需要识别和数据库里的名字有没有重名
	var data model.Category
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBind(&data)
	code = model.CheckCategory(data.Name)
	if code == errmsg.SUCCESS {
		model.EditCate(id, &data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		c.Abort()
	}
}

//删除文章

func DeleteCate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code = model.DeleteCate(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
