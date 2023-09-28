package cleanhttp

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// create http client and return *CleanHttp. Take *Config as params.
func NewFastCleanHttpClient(config *Config) (*FastCleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	client := &fasthttp.Client{
		ReadTimeout:                   config.ReadTimeout,
		WriteTimeout:                  config.WriteTimeout,
		MaxIdleConnDuration:           config.MaxIdleConnDuration,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}

	if config.Proxy != "" {
		config.Proxy = strings.ReplaceAll(config.Proxy, "http://", "")
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(config.Proxy)
	} else {
		client.Dial = (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial
	}

	c := FastCleanHttp{
		Config: config,
		Client: client,
		Log:    config.Log,
	}

	c.BaseHeader = c.GenerateBaseHeaders()

	return &c, nil
}

func (c *FastCleanHttp) Do(request RequestOption) (*fasthttp.Response, error) {
	if request.Url == "" {
		return nil, errors.New("please provide valid url")
	}

	if request.Header == nil {
		request.Header = c.GetDefaultHeader()
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(request.Url)
	req.Header.SetMethod(request.Method)

	for k, v := range request.Header {
		for _, headerValue := range v {
			req.Header.Add(k, headerValue)
		}
	}

	if request.Body != nil {
		b, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		req.SetBodyRaw(b)
	}

	resp := fasthttp.AcquireResponse()
	err := c.Client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	if err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, err
	}

	return resp, nil
}
