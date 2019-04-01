package helper

import (
	"encoding/pem"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	logs "github.com/sirupsen/logrus"
	"github.com/toolkits/net"
	"gopkg.in/ini.v1"
	"reflect"
	"strings"
)

const (
	// system level
	AppName = "appname"
	AppPort = "httpport"
	// ca relative
	KeyType   = "keyType"
	CertType  = "certType"
	KeyStart  = "-----BEGIN RSA PRIVATE KEY-----"
	KeyEnd    = "-----END RSA PRIVATE KEY-----"
	CertStart = "-----BEGIN CERTIFICATE-----"
	CertEnd   = "-----END CERTIFICATE-----"
	// config path
	ConfPath = "./conf/app.conf"
)

type BeegoConfig struct {
	AppName  string `ini:"appname"`
	HttpPort string `ini:"httpport"`
	RunMode  string `ini:"runmode"`
}

type Utils struct {
}

func (u Utils) ReadConfig(path string) (BeegoConfig, error) {
	var config BeegoConfig
	conf, err := ini.Load(path)
	if err != nil {
		logs.Info("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		logs.Info("mapto config file fail!")
		return config, err
	}
	return config, nil
}

func (u Utils) GetAppName() string {
	conf, err := u.ReadConfig(ConfPath)
	if err != nil {
		logs.Error("read beego app.conf failed! ", err)
		return ""
	}

	return conf.AppName
}

func (u Utils) GetAppPort() string {
	conf, err := u.ReadConfig(ConfPath)
	if err != nil {
		logs.Error("read beego app.conf failed! ", err)
		return ""
	}

	return conf.HttpPort
}

func (u Utils) GetRunMode() string {
	conf, err := u.ReadConfig(ConfPath)
	if err != nil {
		logs.Error("read beego app.conf failed! ", err)
		return ""
	}

	return conf.RunMode
}

func (u Utils) SystemHealth() error {
	exception := make([]string, 0)

	// db
	//if ! models.DatabaseHealth() {
	//	exception = append(exception,
	//		"database connect exception; ")
	//}

	// consul
	consul := Consul{}.New()
	if !consul.Health() {
		exception = append(exception, "consul connect exception; ")
	}

	if len(exception) != 0 {
		return fmt.Errorf("%v", exception)
	}

	return nil
}

func (u Utils) NatsHealth(url string) error {
	nc, err := nats.Connect(url, nats.MaxReconnects(-1))
	if err != nil {
		logs.Error(err)
		return err
	}
	nc.Close()
	return nil
}

func (u Utils) Decode(data, keyType string) (result string, err error) {
	if strings.Contains(data, KeyStart) || strings.Contains(data, CertStart) {
		result = data
		return
	}

	switch keyType {
	case KeyType:
		if strings.Contains(data, KeyStart) {
			return
		}
		data = fmt.Sprintf("%s\n%s\n%s", KeyStart, data, KeyEnd)
	case CertType:
		if strings.Contains(data, CertStart) {
			return
		}
		data = fmt.Sprintf("%s\n%s\n%s", CertStart, data, CertEnd)
	}
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		msg := "decode failed"
		logs.Error(msg)
		err = errors.New(msg)
		return
	}
	result = string(block.Bytes)
	return
}

func (u Utils) GetMyIPAddr() string {
	ips, err := net.IntranetIP()
	if err != nil {
		logs.Error("can't get ip list!", err.Error())
		return ""
	}
	return ips[0]
}

func (u Utils) GetMyIdentity() string {
	return strings.Replace(u.GetMyIPAddr(), ".", "_", -1)
}

func (u Utils) ObjectName(v interface{}) (name string) {
	name = u.GetType(v).Name()
	return
}

func (u Utils) ObjectNameAppend(v interface{}, append string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(u.ObjectName(v)), append)
}

func (u Utils) GetType(v interface{}) (ty reflect.Type) {
	ty = reflect.TypeOf(v)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
		return
	}
	return
}
