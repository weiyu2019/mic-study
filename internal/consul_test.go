package internal

import (
	"testing"
)

func TestRegisterService(t *testing.T) {
	err := RegisterService(AppConf.AccountWebConfig.Host,
		AppConf.AccountWebConfig.SrvName,
		AppConf.AccountWebConfig.SrvName,
		AppConf.AccountWebConfig.Port,
		AppConf.AccountWebConfig.Tags)
	if err != nil {
		panic(err)
	}
}
