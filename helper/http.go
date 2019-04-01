package helper

import (
	"context"
	"fmt"

	logs "github.com/sirupsen/logrus"
)

type HttpClient struct {
	Service string
}

func (c HttpClient) Get(respObj interface{}, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "GET", uri, map[string]string{}, nil, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) GetWithHeader(respObj interface{}, url string, header map[string]string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "GET", uri, header, nil, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) Post(reqObj interface{}, respObj interface{}, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "POST", uri, map[string]string{}, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) PostWithHeader(reqObj interface{}, respObj interface{}, header map[string]string, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "POST", uri, header, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) Put(reqObj interface{}, respObj interface{}, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "PUT", uri, map[string]string{}, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) PutWithHeader(reqObj interface{}, respObj interface{}, header map[string]string, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "PUT", uri, header, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) Delete(respObj interface{}, url string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, "DELETE", uri, map[string]string{}, nil, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) Execute(reqObj interface{}, respObj interface{}, method, url string, header map[string]string, urlParams []string) error {
	uri, err := c.buildURI(url, urlParams)
	if err != nil {
		logs.Error(err)
		return err
	}

	if err = c.runE(c.Service, method, uri, header, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) runE(service, method, rawUrl string, headers map[string]string, reqObj interface{}, respObj interface{}) error {
	client := Consul{}.New().GetGoKitConsulClient()

	ep := *GetEndpointMap().getEndpoint(&client, service)
	if _, err := ep(context.Background(), method, rawUrl, headers, reqObj, respObj); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

func (c HttpClient) buildURI(url string, params []string) (string, error) {
	var args []interface{}

	for _, param := range params {
		args = append(args, param)
	}

	return fmt.Sprintf(url, args...), nil
}
