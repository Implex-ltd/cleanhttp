package cleanhttp

import (
	"github.com/Implex-ltd/cleanhttp/internal/cyclepls"
	fp "github.com/Implex-ltd/fingerprint-client/fpclient"
	http "github.com/bogdanfinn/fhttp"
)

type CleanHttp struct {
	Config     *Config
	TlsClient  cyclepls.CyclePLS
	Cookies    []*cyclepls.Cookie
	BaseHeader *HeaderBuilder
}

type Config struct {
	Proxy     string
	Timeout   int
	BrowserFp *fp.Fingerprint
	TlsFp     *fp.TlsFingerprint
}

type RequestOption struct {
	Ja3                    string
	Method                 string
	Body                   []byte
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
	Cookies         string
	UaInfo          UserAgentInfo
}
