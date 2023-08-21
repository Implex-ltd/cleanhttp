package cleanhttp

import (
	"fmt"
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

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(tls_client.Chrome_112), //(GetTlsProfile()),
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

	req.Header.Add("cookie", c.FormatCookies())

	fmt.Println(req.Header)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if _, ok := resp.Header["Set-Cookie"]; ok {
		setCookieHeader := resp.Header["Set-Cookie"][0]
		cookies := c.ParseSetCookieHeader(setCookieHeader)

		for _, cookie := range cookies {
			c.Cookies = append(c.Cookies, &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}

	return resp, nil
}

func (c *CleanHttp) ParseSetCookieHeader(header string) []http.Cookie {
	rawCookies := strings.Split(header, ",")
	var cookies []http.Cookie

	for _, rawCookie := range rawCookies {
		cookieParts := strings.Split(rawCookie, ";")
		cookieNameValue := strings.SplitN(cookieParts[0], "=", 2)

		if len(cookieNameValue) == 2 {
			a, _ := url.QueryUnescape(cookieNameValue[0])
			name := strings.ReplaceAll(strings.TrimSpace(a), "/", "")

			b, _ := url.QueryUnescape(cookieNameValue[1])
			value := strings.TrimSpace(b)

			if name != "" && value != "" {
				cookie := http.Cookie{
					Name:  name,
					Value: value,
				}
				cookies = append(cookies, cookie)
			}
		}
	}
	return cookies
}

// FormatCookies takes all cookies from the client and returns them as a header format string.
func (c *CleanHttp) FormatCookies() string {
	// Default lib use cookiejar..

	var builder strings.Builder

	for i, cookie := range c.Cookies {
		builder.WriteString(cookie.Name)
		builder.WriteString("=")
		builder.WriteString(cookie.Value)

		if i != len(c.Cookies)-1 {
			builder.WriteString("; ")
		}
	}

	return builder.String()
}
