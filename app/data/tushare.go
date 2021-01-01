package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99MyCql/duffett/pkg"
)

const (
	tushareApi = "http://api.waditu.com"
)

// TushareReq 请求 tushare 接口需要的参数
type TushareReq struct {
	ApiName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields  string                 `json:"fields"`
}

// ReqTushareApi 请求 tushare 接口（自动添加 token ）
func ReqTushareApi(data TushareReq) (map[string]interface{}, error) {
	// 请求 tushare 接口时，需携带 token
	data.Token = pkg.Conf.TushareToken
	dataByte, err := json.Marshal(&data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rsp, err := http.Post(tushareApi, "application/json", bytes.NewReader(dataByte))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	res := make(map[string]interface{})
	if err := json.Unmarshal(body, &res); err != nil {
		log.Error(err)
		return nil, err
	}

	code := res["code"].(float64)
	msg := res["msg"].(string)
	if code != 0 {
		log.Error(errors.New(msg))
		return nil, errors.New(msg)
	}

	return res, nil
}

// GetStockName 通过 tsCode 获取股票名
func GetStockName(tsCode string) (string, error) {
	rsp, err := ReqTushareApi(TushareReq{
		ApiName: "stock_basic",
		Params:  map[string]interface{}{"ts_code": tsCode},
		Fields:  "name",
	})
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Debug(rsp)

	item := rsp["data"].(map[string]interface{})["items"].([]interface{})[0].([]interface{})
	if len(item) == 0 {
		log.Error("item 长度为0")
		return "", errors.New("不存在的 ts_code")
	}
	return item[0].(string), nil
}

type DailyData struct {
	TsCode    string
	TradeDate string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	PreClose  float64
	Change    float64 // 涨跌额
	PctChg    float64 // 涨跌幅
	Vol       float64 // 成交量 （手）
	Amount    float64 // 成交额 （千元）
}

// GetDailyData 获取日线数据
func GetDailyData(tsCode string, tradeDate time.Time) (*DailyData, error) {
	rsp, err := ReqTushareApi(TushareReq{
		ApiName: "daily",
		Params: map[string]interface{}{
			"ts_code":    tsCode,
			"trade_date": tradeDate.Format("20060102"),
		},
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	items := rsp["data"].(map[string]interface{})["items"].([]interface{})
	if len(items) == 0 {
		log.Error("items 长度为0")
		return nil, errors.New("未获取到数据")
	}
	item := items[0].([]interface{})

	return &DailyData{
		TsCode:    item[0].(string),
		TradeDate: item[1].(string),
		Open:      item[2].(float64),
		High:      item[3].(float64),
		Low:       item[4].(float64),
		Close:     item[5].(float64),
		PreClose:  item[6].(float64),
		Change:    item[7].(float64),
		PctChg:    item[8].(float64),
		Vol:       item[9].(float64),
		Amount:    item[10].(float64),
	}, nil
}
