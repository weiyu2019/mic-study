package internal

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mic-study/account_srv/model"
	"os"
	"time"
)

var DB *gorm.DB
var err error

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbName"`
	Username string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
}

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	conn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", ViperConf.DBConfig.Username, ViperConf.DBConfig.Password,
		ViperConf.DBConfig.Host, ViperConf.DBConfig.Port, ViperConf.DBConfig.DBName)
	DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger:         newLogger,
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic("数据库链接错误" + err.Error())
	}
	err = DB.AutoMigrate(&model.Account{})
	if err != nil {
		fmt.Println(err)
	}
}
