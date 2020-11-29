package pkg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	Addr         string `yaml:"addr"`
	MysqlUrl     string `yaml:"mysqlUrl"`
	LogPath      string `yaml:"logPath"`
	JwtSecret    string `yaml:"jwtSecret"`
	TushareToken string `yaml:"tushareToken"`
}

// Conf 配置数据
var Conf *config

// InitConfig 读取配置文件，获取配置数据
func InitConfig() {
	// 解析 conf.yaml 文件
	inFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	Conf = new(config)
	err = yaml.Unmarshal(inFile, Conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("config: %+v", *Conf)
}
