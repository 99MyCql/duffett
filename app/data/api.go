package data

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/99MyCql/duffett/pkg"
)

// @Summary Tushare
// @Tags Data
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param data body TushareReq true "data"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/data/tushare [post]
func Tushare(c *gin.Context) {
	var req TushareReq
	if err := c.ShouldBind(&req); err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Print(req)

	rsp, err := ReqTushareApi(req)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端请求 tushare 接口时发生了一些错误"))
		return
	}
	c.JSON(http.StatusOK, rsp)
}
