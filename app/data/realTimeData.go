package data

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const sinaApi = "http://hq.sinajs.cn/list="

// RealTimeData 股票实时数据
type RealTimeData struct {
	StockCode         string
	StockName         string
	TodayOpeningPrice float64
	YdayClosingPrice  float64
	CurPrice          float64
	TodayHighestPrice float64
	TodayLowestPrice  float64
	TradedShareNum    int64
	Turnover          float64
	Date              string
	Time              string
}

// GetRealTimeData 获取股票实时数据
func GetRealTimeData(tsCode string) (*RealTimeData, error) {
	// 将 tushare 格式的 tsCode(000001.SZ) 替换成 sinaApi(sz000001) 接受的格式
	var stockCode string
	i := strings.Index(tsCode, ".")
	if string(tsCode[len(tsCode)-1]) == "Z" {
		stockCode = "sz" + string([]byte(tsCode)[0:i])
	} else if string(tsCode[len(tsCode)-1]) == "H" {
		stockCode = "sh" + string([]byte(tsCode)[0:i])
	} else {
		log.Print("未匹配的 tsCode")
		return nil, errors.New("无法处理的 tsCode")
	}

	// 请求接口获取数据
	rsp, err := http.Get(sinaApi + stockCode)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// GBK转UTF-8
	body, _ = simplifiedchinese.GBK.NewDecoder().Bytes(body)
	strData := string(body)
	log.Print(strData)

	// 从字符串中提取数据
	r, _ := regexp.Compile("\".*?\"")
	strData = strings.Trim(r.FindString(strData), "\"")
	arrData := strings.Split(strData, ",")
	log.Print(arrData)

	todayOpeningPrice, _ := strconv.ParseFloat(arrData[1], 64)
	ydayClosingPrice, _ := strconv.ParseFloat(arrData[2], 64)
	curPrice, _ := strconv.ParseFloat(arrData[3], 64)
	todayHighestPrice, _ := strconv.ParseFloat(arrData[4], 64)
	todayLowestPrice, _ := strconv.ParseFloat(arrData[5], 64)
	tradedShareNum, _ := strconv.ParseInt(arrData[8], 10, 64)
	turnover, _ := strconv.ParseFloat(arrData[9], 64)
	return &RealTimeData{
		StockCode:         stockCode,
		StockName:         arrData[0],
		TodayOpeningPrice: todayOpeningPrice,
		YdayClosingPrice:  ydayClosingPrice,
		CurPrice:          curPrice,
		TodayHighestPrice: todayHighestPrice,
		TodayLowestPrice:  todayLowestPrice,
		TradedShareNum:    tradedShareNum,
		Turnover:          turnover,
		Date:              arrData[30],
		Time:              arrData[31],
	}, nil
}
