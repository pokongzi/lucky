package mysql

import (
	"fmt"
	"lucky/common/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	if DB != nil {
		return
	}
	session := config.Config.Section("mysql")
	user := session.Key("user").String()
	password := session.Key("password").String()
	host := session.Key("host").String()
	port := session.Key("port").String()
	dataBase := session.Key("db").String()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dataBase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	return
}
