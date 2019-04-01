package helper

import (
	"fmt"
	"os"
)

type EnvVar struct {
	ConsulAddr string
	ConsulPort string
}

func (e *EnvVar) load() *EnvVar {
	e.ConsulAddr = os.Getenv("CONSUL_ADDR")
	e.ConsulPort = os.Getenv("CONSUL_PORT")

	return e
}

func (e EnvVar) GetConsulAddr() string {
	p := &e
	p.load()
	return p.ConsulAddr
}

func (e EnvVar) GetConsulPort() string {
	p := &e
	p.load()
	return p.ConsulPort
}

func (e EnvVar) GetConsulURI() string {
	p := &e
	p.load()

	if p.ConsulAddr == "" || p.ConsulPort == "" {
		msg := "invalid env variable CONSUL_PORT or CONSUL_URL"
		panic(msg)
	}

	return fmt.Sprintf("%s:%s", p.ConsulAddr, p.ConsulPort)
}
