package internal

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var AppConf AppConfig
var NacosConf NacosConfig

// var ViperConf ViperConfig
var fileName = "./dev-config.yaml"

func initNacos() {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.ReadInConfig()
	v.Unmarshal(&NacosConf)
	fmt.Println(&NacosConf)
}

func initFromNacos() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NacosConf.Host,
			Port:   NacosConf.Port,
		},
	}
	clientConfigs := constant.ClientConfig{
		NamespaceId:         NacosConf.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfigs": clientConfigs,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: NacosConf.Dataid,
		Group:  NacosConf.Group,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	json.Unmarshal([]byte(content), &AppConf)
}

func init() {
	initNacos()
	initFromNacos()
	fmt.Println("初始化完成")
	InitRedis()
}

type ViperConfig struct {
	DBConfig         DBConfig         `mapstructure:"db"`
	RedisConfig      RedisConfig      `mapstructure:"redis"`
	ConsulConfig     ConsulConfig     `mapstructure:"consul"`
	AccountSrvConfig AccountSrvConfig `mapstructure:"account_srv"`
	AccountWebConfig AccountWebConfig `mapstructure:"account_web"`
	NacosConfig      NacosConfig      `mapstructure:"nacos"`
}
