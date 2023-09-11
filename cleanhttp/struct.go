package cleanhttp

import (
	"io"

	fp "github.com/Implex-ltd/fingerprint-client/fpclient"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type CleanHttp struct {
	Config     *Config
	Client     tls_client.HttpClient
	BaseHeader *HeaderBuilder
	Log        bool
}

type Config struct {
	Proxy     string
	Timeout   int
	Log       bool
	BrowserFp *fp.Fingerprint
	TlsFp     *fp.TlsFingerprint

	Profil *tls_client.ClientProfile
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
