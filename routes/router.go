package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	//创建一个gin实例
	r := gin.New()
	//日志中间件
	r.Use(middleware.Loggoer())
	//使用Gin的内置恢复中间件，以防止程序崩溃
	r.Use(gin.Recovery())
	//auth路由组包含了需要JWT认证的路由，而router路由组包含了不需要认证的路由。
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		//编辑用户
		auth.PUT("users/:id", v1.EditUser)
		///删除用户
		auth.DELETE("users/:id", v1.DeleteUser)
		//分类模块的路由接口
		auth.POST("/category/add", v1.AddCategory)
		//编辑分类名
		auth.PUT("category/:id", v1.EditCate)
		//删除文章
		auth.DELETE("category/:id", v1.DeleteCate)
		//添加文章
		auth.POST("/article/add", v1.AddArticle)
		//查询文章列表
		auth.GET("article/", v1.GetArt)
		//查询所有文章
		auth.DELETE("article/:id", v1.DeleteArt)
		//上传文件
		auth.POST("upload", v1.UploadFile)
	}
	router := r.Group("api/v1")
	{
		//添加用户
		router.POST("/user/add", v1.AddUser)
		//查询用户列表
		router.GET("users/", v1.GetUsers)
		//查询分页列表
		router.GET("category/", v1.GetCate)
		//查询分类下的所有文章
		router.GET("article/list/:id", v1.GetCateArt)
		//单个文章查询
		router.GET("article/:id", v1.GetArtInfo)
		//编辑文章
		router.PUT("article/:id", v1.EditArt)
		//登录
		router.POST("login", v1.Login)
	}
	//启动HTTP服务器，监听utils.HttpPort指定的端口
	r.Run(utils.HttpPort)
}
