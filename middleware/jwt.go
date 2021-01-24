package middleware

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"duffett/pkg"
)

// JWTAuth JWT 验证中间件，返回函数
func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// token 存放在 HTTP 头部的 Authorization 中，形如：Authorization: Bearer <token>
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, pkg.ClientErr("HTTP Header 中 Authorization 字段为空"))
			c.Abort()
			return
		}

		// 按空格分割为 'Bearer' 和 '<token>' 两部分
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, pkg.ClientErr("HTTP Header 中 Authorization 字段格式错误"))
			c.Abort()
			return
		}

		// 解析 token
		myClaims, err := pkg.ParseToken(parts[1])
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusOK, pkg.ClientErr("无效的 token，请重新登录"))
			c.Abort()
			return
		}

		// 判断 token 是否过期
		if myClaims.ExpiresAt <= time.Now().Unix() {
			c.JSON(http.StatusOK, pkg.ClientErr("token 已过期，请重新登录"))
			c.Abort()
			return
		}

		c.Set("username", myClaims.Username)
		c.Next()
	}
}
