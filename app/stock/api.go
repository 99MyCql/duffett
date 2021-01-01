package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/99MyCql/duffett/app/stock/model"
	"github.com/99MyCql/duffett/pkg"
)

// @Summary GetMonitoringStocks
// @Tags Stock
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/stock/getMonitoringStocks [get]
func GetMonitoringStocks(c *gin.Context) {
	username, _ := c.Get("username")
	stocks := model.ListMonitoringStocks(username.(string))
	c.JSON(http.StatusOK, pkg.SucWithData("", stocks))
	return
}

// @Summary GetStocks
// @Tags Stock
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/stock/getStocks [get]
func GetStocks(c *gin.Context) {
	username, _ := c.Get("username")
	stocks := model.ListStocks(username.(string))
	c.JSON(http.StatusOK, pkg.SucWithData("", stocks))
	return
}
