package config

import (
	"log"

	"gopkg.in/ini.v1"
)

var Config *ini.File

func init() {
	ConfigFile, err := ini.Load("config/ini.ini")
	if err != nil {
		log.Fatal("未找到配置文件", err)
		return
	}
	Config = ConfigFile
}
