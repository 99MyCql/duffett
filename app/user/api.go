package user

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/99MyCql/duffett/app/user/model"
	"github.com/99MyCql/duffett/pkg"
)

type loginReq struct {
	Username string `json:"username" binding:"required,excludes= "`
	Password string `json:"password" binding:"required,excludes= "`
}

// @Summary Login
// @Tags User
// @Accept json
// @Param user body loginReq true "user"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	// 定义请求数据结构
	var req loginReq
	// 解析请求数据
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Debug(req)

	// 查找数据库判断是否正确
	if rsp := model.Check(req.Username, req.Password); rsp.Code != pkg.SucCode {
		c.JSON(http.StatusOK, rsp)
		return
	}

	// 返回 token
	token, err := pkg.GenToken(req.Username)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端生成 token 出错"))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("Welcome to duffett!", gin.H{"token": token}))
}

// @Summary TestJwt
// @Tags User
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/user/testJwt [get]
func TestJwt(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, pkg.SucWithData("test", username))
}

type registerReqData struct {
	Username string `json:"username" binding:"required,excludes= "`
	Password string `json:"password" binding:"required,excludes= ,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
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
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Debug(req)

	// 查询用户名是否存在
	if user := model.FindByName(req.Username); user != nil {
		c.JSON(http.StatusOK, pkg.ClientErr("用户名已存在"))
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: pkg.Md5Encode(req.Password),
		Email:    req.Email,
		Sex:      2,
		Role:     "normal",
	}
	if rsp := model.Create(&user); rsp.Code != pkg.SucCode {
		c.JSON(http.StatusOK, rsp)
		return
	}
	log.Debug(user)

	c.JSON(http.StatusOK, pkg.Suc("注册成功"))
}
