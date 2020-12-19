package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"duffett/app/stock/model"
	"duffett/pkg"
)

// @Summary GetMonitoringStocks
// @Tags Stock
// @Accept json
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/stock/getMonitoringStocks [get]
func GetMonitoringStocks(c *gin.Context) {
	username, _ := c.Get("username")
	stocks := model.FindMonitoringStocks(username.(string))
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
	stocks := model.FindStocks(username.(string))
	c.JSON(http.StatusOK, pkg.SucWithData("", stocks))
	return
}
