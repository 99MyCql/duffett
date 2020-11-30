package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"duffett/pkg"
)

const (
	tushareApi = "http://api.waditu.com"
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
		log.Print(err.Error())
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	log.Print(req)

	rsp, err := reqTushareApi(req)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusOK, pkg.ServerErr("服务端请求 tushare 接口时发生了一些错误"))
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func reqTushareApi(data tushareReq) (map[string]interface{}, error) {
	// 请求 tushare 接口时，需携带 token
	data.Token = pkg.Conf.TushareToken
	dataByte, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	rsp, err := http.Post(tushareApi, "application/json", bytes.NewReader(dataByte))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}
