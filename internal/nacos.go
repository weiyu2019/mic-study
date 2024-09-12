package internal

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
