package cleanhttp

import (
	"fmt"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/mssola/user_agent"
)

func GenerateAcceptLanguageHeader(languages []string) string {
	var headerContent strings.Builder

	for i, lang := range languages {
		if i > 0 {
			headerContent.WriteString(",")
		}
		headerContent.WriteString(lang)

		if i > 0 {
			quality := float64(1.0) - float64(i)*0.1
			headerContent.WriteString(";q=")
			headerContent.WriteString(fmt.Sprintf("%.1f", quality))
		}
	}

	return headerContent.String()
}

// Parse user-agent informations and return *UserAgentInfo.
func ParseUserAgent(userAgentString string) *UserAgentInfo {
	ua := user_agent.New(userAgentString)
	name, version := ua.Browser()

	info := &UserAgentInfo{
		BrowserName:    name,
		BrowserVersion: version,
		OSName:         ua.OSInfo().Name,
		OSVersion:      ua.OSInfo().Version,
		UaVersion:      strings.Split(version, ".")[0],
	}
	return info
}

func (c *CleanHttp) GenerateBaseHeaders() *HeaderBuilder {
	ua := ParseUserAgent(c.Config.BrowserFp.Navigator.UserAgent)

	platform := ""
	ch := "?0"
	if strings.Contains(c.Config.BrowserFp.Navigator.UserAgent, "Windows") {
		platform = "Windows"
	} else if strings.Contains(c.Config.BrowserFp.Navigator.UserAgent, "Macintosh") {
		platform = "macOS"
	} else {
		if strings.Contains(c.Config.BrowserFp.Navigator.UserAgent, "Android") {
			platform = "Android"
			ch = "?1"
		} else {
			platform = "Linux"
		}
	}

	h := &HeaderBuilder{
		SecChUa:         fmt.Sprintf(`"Not.A/Brand";v="24", "Chromium";v="%s", "Google Chrome";v="%s"`, ua.UaVersion, ua.UaVersion),
		SecChUaPlatform: fmt.Sprintf(`"%s"`, platform),
		SecChUaMobile:   ch, // todo -> c.Config.BrowserFp.Navigator.Platform,
		AcceptLanguage:  GenerateAcceptLanguageHeader(c.Config.BrowserFp.Navigator.Languages),
		UaInfo:          *ua,
	}

	return h
}

func (c *CleanHttp) GetDefaultHeader() http.Header {
	return http.Header{
		"sec-ch-ua":          {c.BaseHeader.SecChUa},
		"sec-ch-ua-mobile":   {c.BaseHeader.SecChUaMobile},
		"sec-ch-ua-platform": {c.BaseHeader.SecChUaPlatform},
		"user-agent":         {c.Config.BrowserFp.Navigator.UserAgent},
		"sec-fetch-site":     {`none`},
		"sec-fetch-mode":     {`navigate`},
		"sec-fetch-user":     {`?0`},
		"sec-fetch-dest":     {`document`},
		"accept-encoding":    {`gzip, deflate, br`},
		"accept-language":    {c.BaseHeader.AcceptLanguage},

		http.HeaderOrderKey: {
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
			"cookie",
		},
	}
}
