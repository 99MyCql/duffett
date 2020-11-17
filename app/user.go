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

type loginReqData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Login
// @Tags User
// @Accept json
// @Param user body loginReqData true "user"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	// 定义请求数据结构
	var req loginReqData
	// 解析请求数据
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ClientErrCode,
			Msg:  "should post with username and password",
			Data: err.Error(),
		})
		return
	}
	log.Print(req)

	// 查找数据库判断是否正确
	user := model.User{}
	result := pkg.DB.Where("username = ? and password = ?",
		req.Username, pkg.Md5Encode(req.Password)).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ClientErrCode,
			Msg:  "username or password error",
		})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ServerErrCode,
			Msg:  "find data error",
			Data: result.Error.Error(),
		})
		return
	}

	// 返回 token
	token, err := pkg.GenToken(req.Username)
	if err != nil {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ServerErrCode,
			Msg:  "generate token fail",
			Data: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, pkg.RspData{
		Code: pkg.SucCode,
		Msg:  "Welcome to duffett!",
		Data: map[string]interface{}{
			"token": token,
		},
	})
}

// @Summary Test
// @Tags User
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/user/test [get]
func Test(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, pkg.RspData{
		Code: pkg.SucCode,
		Msg:  "test",
		Data: username,
	})
}

type registerReqData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Register
// @Tags User
// @Accept json
// @Param user body registerReqData true "user"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/user/register [post]
func Register(c *gin.Context) {
	var req registerReqData
	// 解析请求数据
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ClientErrCode,
			Msg:  "should post with username and password",
			Data: err.Error(),
		})
		return
	}
	log.Print(req)

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ClientErrCode,
			Msg:  "username or password is empty",
		})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: pkg.Md5Encode(req.Password),
		Sex:      2,
		Role:     "normal",
	}
	result := pkg.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, pkg.RspData{
			Code: pkg.ServerErrCode,
			Msg:  "insert data error",
			Data: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pkg.RspData{
		Code: pkg.SucCode,
		Msg:  "register",
	})
}
