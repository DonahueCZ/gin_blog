package middleware

import (
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// 定一个全局变量JwyKey，它是从utils包中的JwtKey转换成字符切片

var JwyKey = []byte(utils.JwtKey)

type MyClaims struct {
	//要和user里面的保持一致
	Username string `json:"username"`
	//RegisteredClaims结构体里有一些基础定义的东西
	jwt.RegisteredClaims
}

// 生成token

func SetToken(Username string) (string, int) {
	//token过期时间设定
	expireTime := time.Now().Add(10 * time.Hour)
	SetClaims := MyClaims{
		Username: Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), //设置过期时间0
			Issuer:    "ginBlog",                      //设置发行人
		},
	}
	//使用HS256签名方法和SetClaims创建一个新的token
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	// 使用JwtKey对token进行签名
	token, err := reqClaim.SignedString([]byte(utils.JwtKey))
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

//	验证token

func CheckToken(token string) (*MyClaims, int) {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwyKey, nil
	})
	if err != nil {
		// 如果有错误，比如令牌过期或令牌无效，应该返回相应的错误码
		return nil, errmsg.ERROR_TOKEN_WRONG
	}
	if key, code := setToken.Claims.(*MyClaims); code && setToken.Valid {
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR_TOKEN_WRONG
	}
}

//JwtToken函数返回一个Gin的中间件，用于处理JWT认证

func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization") //从请求头中获取Authorization字段
		code := errmsg.SUCCESS

		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		checkToken := strings.SplitN(tokenHeader, " ", 2) // 将Authorization字段分割为两部分
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		key, tCode := CheckToken(checkToken[1]) //验证JWT
		if tCode != errmsg.SUCCESS {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		//如果用户名是空的
		if key == nil || key.Username == "" {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "令牌无效或用户名为空",
			})
			c.Abort()
			return
		}

		//如果JWT过期，返回错误
		if time.Now().Unix() > key.ExpiresAt.Unix() {
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "JWT已过期",
			})
			c.Abort()
			return
		}

		c.Set("username", key.Username)
		c.Next()
	}
}
