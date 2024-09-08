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
	RedisConfig RedisConfig `mapstructure:"redis"`
}
