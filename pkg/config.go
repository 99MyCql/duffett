package pkg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type config struct {
	Addr     string `yaml:"addr"`
	MysqlUrl string `yaml:"mysqlUrl"`
	LogPath  string `yaml:"logPath"`
}

var Conf *config

func InitConfig() {
	// 解析 conf.yaml 文件
	inFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	Conf = new(config)
	err = yaml.Unmarshal(inFile, Conf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("config: %+v", *Conf)
}
