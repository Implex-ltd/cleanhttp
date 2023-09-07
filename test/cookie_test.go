package cleanhttp

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"github.com/Implex-ltd/fingerprint-client/fpclient"
)

func TestCookie(t *testing.T) {
	fp, err := fpclient.LoadFingerprint(&fpclient.LoadingConfig{
		FilePath: "../assets/chrome114.json",
	})

	if err != nil {
		panic(err)
	}

	cfg := &cleanhttp.Config{
		BrowserFp: fp,
	}

	type args struct {
		config *cleanhttp.Config
		url    string
		method string
		body   io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *cleanhttp.CleanHttp
		wantErr bool
	}{
		{
			name: "get cookies",
			args: args{
				config: cfg,
				method: "GET",
				body:   strings.NewReader(``),
				url:    "https://discord.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := cleanhttp.NewCleanHttpClient(tt.args.config)

			if err != nil {
				panic(err)
			}

			resp, err := c.Do(cleanhttp.RequestOption{
				Method: tt.args.method,
				Url:    tt.args.url,
				Header: c.GetDefaultHeader(),
				Body:   tt.args.body,
			})

			if err != nil {
				panic(err)
			}

			fmt.Println(resp.Header)
			u, _ := url.Parse("https://discord.com")
			fmt.Println(c.Client.GetCookieJar().Cookies(u))
			fmt.Println("----------------------------")

			r, err := c.Do(cleanhttp.RequestOption{
				Method: "GET",
				Url:    "http://httpbin.org/cookies",
			})

			if err != nil {
				panic(err)
			}

			b, _ := io.ReadAll(r.Body)

			fmt.Println(string(b))
		})
	}
}
