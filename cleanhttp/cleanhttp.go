package cleanhttp

import (
	"bytes"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strings"

	"github.com/Implex-ltd/cleanhttp/internal/cyclepls"
	tls_client "github.com/bogdanfinn/tls-client"
)

// create http client and return *CleanHttp. Take *Config as params.
func NewCleanHttpClient(config *Config) (*CleanHttp, error) {
	if config.Timeout < 30 {
		config.Timeout = 30
	}

	c := CleanHttp{
		Cookies:   nil,
		Config:    config,
		TlsClient: cyclepls.Init(),
	}

	c.BaseHeader = c.GenerateBaseHeaders()

	return &c, nil
}

func (c *CleanHttp) DoTls(request RequestOption) (*cyclepls.Response, error) {
	if request.Header == nil {
		request.Header = c.GetDefaultHeader()
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(GetTlsProfile()),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithProxyUrl(c.Config.Proxy),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(request.Body)

	// Read the request body into a buffer
	var bodyBuffer bytes.Buffer
	_, err = io.Copy(&bodyBuffer, reader)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(request.Method, request.Url, &bodyBuffer)
	if err != nil {
		return nil, err
	}

	for k, v := range request.Header {
		req.Header.Add(k, v[0])
	}

	req.Header.Add("cookie", c.FormatCookies())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if _, ok := resp.Header["Set-Cookie"]; ok {
		setCookieHeader := resp.Header["Set-Cookie"][0]
		cookies := c.ParseSetCookieHeader(setCookieHeader)

		for _, cookie := range cookies {
			c.Cookies = append(c.Cookies, &cyclepls.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}

	defer resp.Body.Close()

	r, _ := io.ReadAll(resp.Body)

	re := cyclepls.Response{
		Body:   string(r),
		Status: resp.StatusCode,
	}

	return &re, nil
}

// Do request and return *http.Response, Take RequestOption in params.
func (c *CleanHttp) Do(request RequestOption) (*cyclepls.Response, error) {
	if request.Header == nil {
		request.Header = c.GetDefaultHeader()
	}

	/* Give http 400 ? */
	/*if request.Body != nil && request.CalculateContentLength {
		len, e := CalculateContentLength(request.Body)

		if e == nil {
			request.Header.Add("content-length", strconv.Itoa(int(len)))
		}
	}*/

	if request.Ja3 == "" {
		request.Ja3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,27-35-65281-17513-23-43-51-18-11-16-45-10-13-5-0-41,25497-29-23-24,0"
	}

	headerMap := make(map[string]string)
	for key, values := range request.Header {
		if len(values) > 0 {
			headerMap[key] = values[0]
		}
	}

	headerMap["cookie"] = c.FormatCookies()

	opt := cyclepls.Options{
		Ja3:       request.Ja3,
		UserAgent: headerMap["user-agent"],
		Headers:   headerMap,
		Cookies:   c.TlsClient.Cookies,
	}

	if (request.Method == "PUT" || request.Method == "POST") && request.Body != nil {
		opt.Body = string(request.Body)
	}

	if c.Config.Proxy != "" {
		opt.Proxy = c.Config.Proxy
	}

	resp, err := c.TlsClient.Do(request.Url, opt, strings.ToUpper(request.Method))
	if err != nil {
		return nil, err
	}

	if _, ok := resp.Headers["Set-Cookie"]; ok {
		setCookieHeader := resp.Headers["Set-Cookie"]
		cookies := c.ParseSetCookieHeader(setCookieHeader)

		for _, cookie := range cookies {
			c.Cookies = append(c.Cookies, &cyclepls.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}

	return &resp, nil
}

func (c *CleanHttp) ParseSetCookieHeader(header string) []cyclepls.Cookie {
	rawCookies := strings.Split(header, ",")
	var cookies []cyclepls.Cookie

	for _, rawCookie := range rawCookies {
		cookieParts := strings.Split(rawCookie, ";")
		cookieNameValue := strings.SplitN(cookieParts[0], "=", 2)

		if len(cookieNameValue) == 2 {
			a, _ := url.QueryUnescape(cookieNameValue[0])
			name := strings.ReplaceAll(strings.TrimSpace(a), "/", "")

			b, _ := url.QueryUnescape(cookieNameValue[1])
			value := strings.TrimSpace(b)

			if name != "" && value != "" {
				cookie := cyclepls.Cookie{
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
