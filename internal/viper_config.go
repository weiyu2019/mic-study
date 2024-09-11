package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

var ViperConf ViperConfig
var fileName = "./dev-config.yaml"

func init() {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.ReadInConfig()
	v.Unmarshal(&ViperConf)
	fmt.Println(&ViperConf)
	InitRedis()
}

type ViperConfig struct {
	DBConfig         DBConfig         `mapstructure:"db"`
	RedisConfig      RedisConfig      `mapstructure:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul"`
	AccountSrvConfig AccountSrvConfig `mapstructure:"account_srv"`
	AccountWebConfig AccountWebConfig `mapstructure:"account_web"`
}
