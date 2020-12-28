package data

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/99MyCql/duffett/pkg"
)

// tushareReq 请求 tushare 接口需要的参数
type tushareReq struct {
	ApiName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields  string                 `json:"fields"`
}

// @Summary Tushare
// @Tags Data
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Param data body tushareReq true "data"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":非0,"data":{},"msg":""}"
// @Router /api/v1/data/tushare [post]
func Tushare(c *gin.Context) {
	var req tushareReq
	if err := c.ShouldBind(&req); err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Print(req)

	rsp, err := reqTushareApi(req)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端请求 tushare 接口时发生了一些错误"))
		return
	}
	c.JSON(http.StatusOK, rsp)
}
