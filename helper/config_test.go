package helper

import (
	logs "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestConfig_LoadConf(t *testing.T) {
	conf := GetConfInst()
	t.Log("conf:", conf.LoadConf())

	conf.Run()
	t.Log("conf:", conf.LoadConf())

	f := func(delta *ConfDelta) {
		logs.Info("old conf:", delta.OData)
		logs.Info("new conf:", delta.NData)
	}
	conf.Register("test-callback", f)

	consul := Consul{}.New()
	key := "config/nautilus/Data"
	v, err := consul.GetKey(key)
	if err != nil {
		t.Error(err)
	}

	consul.SetKey(key, "{\"DBUrl\":\"localhost\",\"DBPort\":\"3307\",\"DBUser\":\"root\", \"DBPwd\":\"\"}")
	time.Sleep(3 * time.Second)

	consul.SetKey(key, v)

	conf.DeRegister("test-callback")
}
