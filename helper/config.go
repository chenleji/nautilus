package helper

import (
	"fmt"
	"github.com/mohae/deepcopy"
	logs "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"reflect"
	"time"
)

const (
	ConfigFileName      = "data"
	ConfigProvider      = "consul"
	ConfigPathFmt       = "/config/%s/%s"
	ConfigFmtYaml       = "yaml"
	ConfigFmtJson       = "json"
	ConfigWatchInterval = 1 * time.Second
)

///////////////////////////////////
// config file
//////////////////////////////////
type ConfDelta struct {
	OData interface{}
	NData interface{}
}

///////////////////////////////////
// config
//////////////////////////////////

var config *Config

func GetConfInst(data ...interface{}) *Config {
	if config == nil {
		if len(data) == 1 {
			config = (&Config{}).init(data[0])
		} else {
			logs.Error("please input config data structure pointer.")
			return nil
		}
	}

	return config
}

type Config struct {
	FileName    string
	Format      string
	Data        interface{}
	callbackMap map[string]func(confDelta *ConfDelta)
}

func (c *Config) init(data interface{}) *Config {
	if c.FileName == "" {
		c.FileName = ConfigFileName
	}

	if c.Format == "" {
		c.Format = ConfigFmtJson
	}

	c.Data = data
	c.callbackMap = make(map[string]func(confDelta *ConfDelta))
	c.getViper()

	return c
}

func (c *Config) LoadConf() interface{} {
	return c.Data
}

func (c *Config) Register(name string, f func(confDelta *ConfDelta)) {
	c.callbackMap[name] = f
}

func (c *Config) DeRegister(name string) {
	if _, ok := c.callbackMap[name]; ok {
		delete(c.callbackMap, name)
	}

}

func (c *Config) broadcast(delta *ConfDelta) {
	for _, f := range c.callbackMap {
		f(delta)
	}
}

func (c *Config) Run() {
	go func() {
		producer := c.getViper()
		lastObj := deepcopy.Copy(c.Data)

		for {
			time.Sleep(ConfigWatchInterval)

			producer.Unmarshal(c.Data)
			if ! reflect.DeepEqual(lastObj, c.Data) {
				delta := &ConfDelta{
					OData: lastObj,
					NData: c.Data,
				}
				c.broadcast(delta)
			}

			lastObj = deepcopy.Copy(c.Data)
		}
	}()

	return
}

func (c *Config) getViper() *viper.Viper {
	runtimeViper := viper.New()
	consulAddr := EnvVar{}.GetConsulURI()
	appName := Utils{}.GetAppName()
	path := fmt.Sprintf(ConfigPathFmt, appName, c.FileName)

	runtimeViper.AddRemoteProvider(ConfigProvider, consulAddr, path)
	runtimeViper.SetConfigType(c.Format)
	err := runtimeViper.ReadRemoteConfig()
	if err != nil {
		logs.Error(err)
	}

	runtimeViper.Unmarshal(c.Data)
	err = runtimeViper.WatchRemoteConfigOnChannel()
	if err != nil {
		logs.Error("unable to read remote config: ", err)
	}

	return runtimeViper
}
