package cleanhttp

import (
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"

	tls_client "github.com/bogdanfinn/tls-client"
)

// create http client and return *CleanHttp. Take *Config as params.
func NewCleanHttpClient(config *Config) (*CleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	if config.Profil == nil {
		config.Profil = GetTlsProfile()
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(*config.Profil),
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
		Cookies: nil,
		Config:  config,
		Client:  client,
	}

	c.BaseHeader = c.GenerateBaseHeaders()

	return &c, nil
}

func (c *CleanHttp) Do(request RequestOption) (*http.Response, error) {
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

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FormatCookies takes all cookies from the client and returns them as a header format string.
func (c *CleanHttp) FormatCookies(url *url.URL) string {
	var builder strings.Builder

	for i, cookie := range c.Client.GetCookieJar().Cookies(url) {
		builder.WriteString(cookie.Name)
		builder.WriteString("=")
		builder.WriteString(cookie.Value)

		if i != len(c.Cookies)-1 {
			builder.WriteString("; ")
		}
	}

	return builder.String()
}
