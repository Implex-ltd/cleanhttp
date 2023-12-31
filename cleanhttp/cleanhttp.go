package cleanhttp

import (
	"errors"
	"log"
	"net/url"

	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"

	tls_client "github.com/bogdanfinn/tls-client"
)

// create http client and return *CleanHttp. Take *Config as params.
func NewCleanHttpClient(config *Config) (*CleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithRandomTLSExtensionOrder(),
	}

	if config.Proxy != "" {
		options = append(options, tls_client.WithProxyUrl(config.Proxy))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	c := CleanHttp{
		Config: config,
		Client: client,
		Log:    config.Log,
	}

	c.BaseHeader = c.GenerateBaseHeaders()

	return &c, nil
}

func (c *CleanHttp) Do(request RequestOption) (*http.Response, error) {
	if request.Url == "" {
		return nil, errors.New("please provide valid url")
	}

	if request.Header == nil {
		request.Header = c.GetDefaultHeader()
	}

	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range request.Header {
		req.Header.Set(k, v[0])
	}

	u, _ := url.Parse(request.Url)
	for _, cook := range c.Client.GetCookieJar().Cookies(u) {
		req.AddCookie(cook)

		if c.Log {
			log.Println("add cookie", cook.Name, cook.Value, cook.Domain)
		}
	}

	if c.Log {
		log.Println(req.URL, req.Header, req.Body)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
