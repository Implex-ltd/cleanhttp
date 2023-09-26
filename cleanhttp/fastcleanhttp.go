package cleanhttp

import (
	"errors"
	"io"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// create http client and return *CleanHttp. Take *Config as params.
func NewFastCleanHttpClient(config *Config) (*FastCleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	readTimeout, _ := time.ParseDuration("5s")
	writeTimeout, _ := time.ParseDuration("5s")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial:                          fasthttpproxy.FasthttpHTTPDialer(config.Proxy),
	}

	c := FastCleanHttp{
		Config: config,
		Client: client,
		Log:    config.Log,
	}

	c.BaseHeader = c.GenerateBaseHeaders()

	return &c, nil
}

func (c *FastCleanHttp) Do(request RequestOption) ([]byte, int, error) {
	if request.Url == "" {
		return nil, 0, errors.New("please provide valid url")
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
			return nil, 0, err
		}

		req.SetBodyRaw(b)
	}

	resp := fasthttp.AcquireResponse()
	err := c.Client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return nil, 0, err
	}

	return resp.Body(), resp.StatusCode(), nil
}
