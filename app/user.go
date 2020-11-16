package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"duffett/model"
	"duffett/pkg"
)

// @Summary Login
// @Tags User
// @version 1.0
// @Accept json
// @Param username query string true "admin"
// @Param password query string true "xxx"
// @Success 200 {string} json "{"code":0,"data":{},"msg":"success"}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":"xxx"}"
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	// 定义请求数据结构
	var reqData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// 解析请求数据
	err := c.ShouldBind(&reqData)
	if err != nil {
		c.JSON(http.StatusOK, pkg.RspMsg{
			Code: pkg.ErrCode,
			Msg:  "error",
			Data: err.Error(),
		})
		return
	}
	log.Println(reqData)

	// 查找数据库判断是否正确
	user := model.User{}
	result := pkg.DB.Where("username = ? and password = ?", reqData.Username, reqData.Password).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, pkg.RspMsg{
			Code: pkg.FailCode,
			Msg:  "username or password error",
		})
	} else {
		c.JSON(http.StatusOK, pkg.RspMsg{
			Code: pkg.SucCode,
			Msg:  "Welcome to the administrator!",
		})
	}
}
