package circuitbreaker

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/chenleji/kit/endpoint"
)

// Hystrix returns an endpoint.Middleware that implements the circuit
// breaker pattern using the afex/hystrix-go package.
//
// When using this circuit breaker, please configure your commands separately.
//
// See https://godoc.org/github.com/afex/hystrix-go/hystrix for more
// information.
func Hystrix(commandName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, method, rawUrl string, headers map[string]string, reqObj interface{}, respObj interface{}) (response interface{}, err error) {
			var resp interface{}
			if err := hystrix.Do(commandName, func() (err error) {
				resp, err = next(ctx, method, rawUrl, headers, reqObj, respObj)
				return err
			}, nil); err != nil {
				return nil, err
			}
			return resp, nil
		}
	}
}
