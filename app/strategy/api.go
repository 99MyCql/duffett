package strategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/99MyCql/duffett/app/strategy/model"
	strategyModel "github.com/99MyCql/duffett/app/strategy/model"
	userModel "github.com/99MyCql/duffett/app/user/model"
	"github.com/99MyCql/duffett/pkg"
)

// @Summary GetStrategies
// @Tags Strategy
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/strategy/getStrategies [get]
func GetStrategies(c *gin.Context) {
	username, _ := c.Get("username")
	strategies := model.ListStrategies(username.(string))
	c.JSON(http.StatusOK, pkg.SucWithData("", strategies))
	return
}

type createReq struct {
	Name    string `json:"name" binding:"required,excludes= "`
	Desc    string `json:"desc" binding:""`
	Content string `json:"content" binding:"required"`
}

// @Summary Create
// @Tags Strategy
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param strategy body createReq true "createReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/strategy/create [post]
func Create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Debug(req)

	username, _ := c.Get("username")
	u := userModel.FindByName(username.(string))
	if u == nil {
		log.Error("查询不到该用户")
		c.JSON(http.StatusOK, pkg.ServerErr("查询不到该用户"))
		return
	}

	// 检查该策略名是否已存在
	if s := strategyModel.FindByName(u.Username + "_" + req.Name); s != nil {
		log.Error("该策略名已存在")
		c.JSON(http.StatusOK, pkg.ClientErr("该策略名已存在"))
		return
	}

	c.JSON(http.StatusOK, strategyModel.Create(&strategyModel.Strategy{
		Name:    u.Username + "_" + req.Name,
		Desc:    req.Desc,
		Content: req.Content,
		UserID:  u.ID,
	}))
}

type updateReq struct {
	Id      uint   `json:"id" binding:"required"`
	Desc    string `json:"desc" binding:""`
	Content string `json:"content" binding:"required"`
}

// @Summary Update
// @Tags Strategy
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param strategy body updateReq true "updateReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/strategy/update [post]
func Update(c *gin.Context) {
	var req updateReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Debug(req)

	s := strategyModel.FindById(req.Id)
	if s == nil {
		c.JSON(http.StatusOK, pkg.ClientErr("策略ID不存在"))
		return
	}

	s.Desc = req.Desc
	s.Content = req.Content
	c.JSON(http.StatusOK, strategyModel.Update(s))
}
