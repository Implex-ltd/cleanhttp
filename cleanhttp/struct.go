package cleanhttp

import (
	"io"
	"time"

	fp "github.com/Implex-ltd/fingerprint-client/fpclient"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/valyala/fasthttp"
)

type CleanHttp struct {
	Config     *Config
	Client     tls_client.HttpClient
	BaseHeader *HeaderBuilder
	Log        bool
}

type FastCleanHttp struct {
	Config     *Config
	Client     *fasthttp.Client
	BaseHeader *HeaderBuilder
	Log        bool
}

type Config struct {
	Proxy     string
	Timeout   int
	Log       bool
	BrowserFp *fp.Fingerprint

	ReadTimeout, WriteTimeout, MaxIdleConnDuration time.Duration
}

type RequestOption struct {
	Ja3                    string
	Method                 string
	Body                   io.Reader
	Url                    string
	Header                 http.Header
	CalculateContentLength bool
}

type UserAgentInfo struct {
	BrowserName    string
	BrowserVersion string
	OSName         string
	OSVersion      string
	UaVersion      string
}

type HeaderBuilder struct {
	SecChUa         string
	SecChUaPlatform string
	SecChUaMobile   string
	AcceptLanguage  string
	UaInfo          UserAgentInfo
}
