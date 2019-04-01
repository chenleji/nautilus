package helper

import (
	"github.com/chenleji/kit/circuitbreaker"
	"github.com/chenleji/kit/endpoint"
	"github.com/chenleji/kit/log"
	"github.com/chenleji/kit/sd"
	"github.com/chenleji/kit/sd/consul"
	"github.com/chenleji/kit/sd/lb"
	httpTransport "github.com/chenleji/kit/transport/http"
	"io"
	"os"
	"strings"
	"time"
)

const (
	retryMax     = 5
	retryTimeout = 1500 * time.Millisecond
)

var (
	epMap *EndpointMap
)

// endpoint map
type EndpointMap struct {
	data map[string]*endpoint.Endpoint
}

func GetEndpointMap() *EndpointMap {
	if epMap == nil {
		epMap = &EndpointMap{
			data: make(map[string]*endpoint.Endpoint),
		}
	}

	return epMap
}

func (m *EndpointMap) getEndpoint(client *consul.Client, service string) *endpoint.Endpoint {
	if inst, ok := m.data[service]; ok {
		return inst
	}
	return m.newEndpoint(client, service)
}

func (m *EndpointMap) newEndpoint(client *consul.Client, service string) *endpoint.Endpoint {
	var (
		logger   = log.NewLogfmtLogger(os.Stderr)
		tags     = make([]string, 0)
		instance = consul.NewInstancer(*client, logger, service, tags, true)
	)

	factory := endpointFactory()()
	endpointer := sd.NewEndpointer(instance, factory, logger)
	loadBalance := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, retryTimeout, loadBalance)
	ret := circuitbreaker.Hystrix(service)(retry)

	m.data[service] = &ret
	return &ret
}

func endpointFactory() (func() sd.Factory) {
	return func() sd.Factory {
		return func(instance string) (endpoint.Endpoint, io.Closer, error) {
			if !strings.HasPrefix(instance, "http") {
				instance = "http://" + instance
			}
			return httpTransport.NewClient(instance).Endpoint(), nil, nil
		}
	}
}
