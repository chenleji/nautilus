package helper

import (
	"fmt"
	"github.com/chenleji/kit/log"
	goKitConsul "github.com/chenleji/kit/sd/consul"
	stdConsul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	logs "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const (
	ConsulDeRegisterCriticalServiceAfter = "10m"
	ConsulInterval                       = "15s"
	ConsulTTL                            = "10s"
	ConsulDefaultSchema                  = "http://"
)

type WatchResp struct {
	Value []byte
	Error error
}

type Consul struct {
	client      *stdConsul.Client
	goKitClient goKitConsul.Client
}

func (c Consul) New() *Consul {
	// consul client
	client, err := stdConsul.NewClient(&stdConsul.Config{
		Address: EnvVar{}.GetConsulURI(),
	})
	if err != nil {
		logs.Error(err)
		return nil
	}

	// goKit consul client
	c.goKitClient = goKitConsul.NewClient(client)
	c.client = client

	return &c
}

func (c *Consul) GetConsulClient() *stdConsul.Client {
	return c.client
}

func (c *Consul) GetGoKitConsulClient() goKitConsul.Client {
	return c.goKitClient

}

func (c *Consul) RegistryService(service, port, healthCheckURL string) error {
	var err error

	registryIp := Utils{}.GetMyIPAddr()
	registryPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	r := &stdConsul.AgentServiceRegistration{
		ID:                service + "_" + registryIp + "_" + port,
		Name:              service,
		Tags:              []string{},
		Port:              registryPort,
		Address:           registryIp,
		EnableTagOverride: false,
		Check: &stdConsul.AgentServiceCheck{
			DeregisterCriticalServiceAfter: ConsulDeRegisterCriticalServiceAfter,
			HTTP:                           fmt.Sprintf("%s%s:%d%s", ConsulDefaultSchema, registryIp, registryPort, healthCheckURL),
			Interval:                       ConsulInterval,
			Timeout:                        ConsulTTL,
		},
	}
	registrar := goKitConsul.NewRegistrar(c.goKitClient, r, log.NewLogfmtLogger(os.Stderr))
	registrar.Register()

	return nil
}

func (c *Consul) Health() bool {
	// get raft leader
	if _, err := c.client.Status().Leader(); err != nil {
		logs.Error(err)
		return false
	}

	return true
}

func (c *Consul) SetKey(key, value string) error {
	keyPair := stdConsul.KVPair{
		Key:   key,
		Value: []byte(value),
	}
	if _, err := c.client.KV().Put(&keyPair, nil); err != nil {
		logs.Error("SetKey stdConsul set key err: ", err)
		return err
	}

	return nil
}

func (c *Consul) GetKey(key string) (string, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		logs.Error("GetKey stdConsul get key err: ", err)
		return "", err
	}

	return string(pair.Value), nil
}

func (c *Consul) DeleteKey(key string) error {
	if _, err := c.client.KV().Delete(key, nil); err != nil {
		logs.Error("SetKey stdConsul set key err: ", err)
		return err
	}

	return nil
}

func (c *Consul) CheckKey(key string, opts *stdConsul.QueryOptions) bool {
	if opts == nil {
		opts = &stdConsul.QueryOptions{}
	}

	keyPair, _, err := c.client.KV().Get(key, opts)
	if err != nil {
		logs.Error("CheckKey stdConsul get key err: ", err)
		return false
	}

	if keyPair == nil && err == nil {
		err := errors.New("CheckKey key was not found.")
		logs.Error(err)
		return false
	}

	return true
}

func (c *Consul) WatchKey(key string, stop chan bool) <-chan *WatchResp {
	var (
		respChan       = make(chan *WatchResp)
		WaitTimeSecond = 1 * time.Minute
		retryPeriod    = 5 * time.Second
	)

	go func() {
		opts := &stdConsul.QueryOptions{WaitTime: WaitTimeSecond}

		originKeyPair, originMeta, err := c.client.KV().Get(key, opts)
		if err != nil {
			logs.Error("watch stdConsul key failed. err: ", err)
			respChan <- &WatchResp{nil, err}
			return
		}
		if originKeyPair == nil && err == nil {
			err := errors.New("key was not found.")
			logs.Error(err)
			respChan <- &WatchResp{nil, err}
			return
		}

		for {
			select {
			case <-stop:
				return

			default:
				opts := &stdConsul.QueryOptions{
					WaitTime:  WaitTimeSecond,
					WaitIndex: originMeta.LastIndex,
				}

				keyPair, meta, err := c.client.KV().Get(key, opts)
				if err != nil {
					respChan <- &WatchResp{nil, err}
					time.Sleep(retryPeriod)
					continue
				}

				if originKeyPair.CreateIndex != keyPair.CreateIndex ||
					originKeyPair.ModifyIndex != keyPair.ModifyIndex ||
					originKeyPair.LockIndex != keyPair.LockIndex {

					respChan <- &WatchResp{keyPair.Value, nil}
				}

				originKeyPair = keyPair
				originMeta = meta
			}
		}
	}()
	return respChan
}
