package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// 这是跨域的请求文件，可以处理从其他服务器获取的数据
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowOrigins:    []string{"*"},
			AllowMethods:    []string{"*"}, //允许的请求方法
			AllowHeaders:    []string{"Origin"},
			ExposeHeaders:   []string{"Content-Length", "Authorization"},
			MaxAge:          12 * time.Hour,
		})
	}
}
