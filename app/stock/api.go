package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"duffett/app/stock/model"
	model2 "duffett/app/user/model"
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
	u := model2.FindByName(username.(string))
	stocks := model.FindMonitoringStocks(u.ID)
	c.JSON(http.StatusOK, pkg.SucWithData("", stocks))
	return
}
