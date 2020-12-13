package data

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"duffett/pkg"
)

const (
	tushareApi = "http://api.waditu.com"
)

// 请求 tushare 接口，自动添加 token
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

// GetStockName 通过 tsCode 获取股票名
func GetStockName(tsCode string) (string, error) {
	rsp, err := reqTushareApi(tushareReq{
		ApiName: "stock_basic",
		Params:  map[string]interface{}{"ts_code": tsCode},
		Fields:  "name",
	})
	if err != nil {
		return "", err
	}
	log.Print(rsp)
	return rsp["data"].(map[string]interface{})["items"].([]interface{})[0].([]interface{})[0].(string), nil
}
