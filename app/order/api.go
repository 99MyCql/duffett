package order

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/99MyCql/duffett/app/order/model"
	"github.com/99MyCql/duffett/pkg"
)

type getOrdersReq struct {
	StockId uint `json:"stockId" binding:"required"`
}

// @Summary GetOrders
// @Tags Order
// @Accept json
// @Param getOrdersReq body getOrdersReq true "getOrdersReq"
// @Param Authorization header string false "Bearer <token>"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Router /api/v1/order/getOrders [post]
func GetOrders(c *gin.Context) {
	username, _ := c.Get("username")
	var req getOrdersReq
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Debug(req)
	orders := model.ListOrders(username.(string), req.StockId)
	c.JSON(http.StatusOK, pkg.SucWithData("", orders))
	return
}
