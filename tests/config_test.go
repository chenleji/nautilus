package test

import (
	_ "github.com/spf13/viper/remote"
	"testing"
)

func TestWatchConfig(t *testing.T) {
	//ch := make(chan *helper.ConfEvent)
	//helper.GetConfEvent().Register(ch)

	//for i := 0; i< 3; i++ {
	//	select {
	//	case e1 := <-ch:
	//		utils.Display("data is changed: ", e1)
	//	}
	//}
}

func TestConfig(t *testing.T) {
	//var runtimeViper = viper.New()
	//
	//runtimeViper.AddRemoteProvider("consul", "127.0.0.1:8500", "/config/nautilus/data/config.json")
	//runtimeViper.SetConfigType("json")
	//
	//runtimeConf := map[string]interface{}{}
	//
	//err := runtimeViper.ReadRemoteConfig()
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//runtimeViper.Unmarshal(&runtimeConf)
	//utils.Display("runtimeConf", runtimeConf)
	//
	//err = runtimeViper.WatchRemoteConfigOnChannel()
	//if err != nil {
	//	logs.Error("unable to read remote config: %v", err)
	//}

	//for i:= 0; i < 30; i++ {
	//	runtimeViper.Unmarshal(&runtimeConf)
	//	utils.Display("runtimeConf", runtimeConf)
	//	time.Sleep(2 * time.Second)
	//}

}
