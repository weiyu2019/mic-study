package internal

import (
	"testing"
)

func TestRegisterService(t *testing.T) {
	err := RegisterService(ViperConf.AccountWebConfig.Host,
		ViperConf.AccountWebConfig.SrvName,
		ViperConf.AccountWebConfig.SrvName,
		ViperConf.AccountWebConfig.Port,
		ViperConf.AccountWebConfig.Tags)
	if err != nil {
		panic(err)
	}
}
