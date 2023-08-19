package cleanhttp

import (
	"fmt"
	"testing"

	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"github.com/Implex-ltd/fingerprint-client/fpclient"
)

func TestNewCleanHttpClient(t *testing.T) {
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
		body   string
	}
	tests := []struct {
		name    string
		args    args
		want    *cleanhttp.CleanHttp
		wantErr bool
	}{
		{
			name: "get tls",
			args: args{
				config: cfg,
				method: "GET",
				body:   "",
				url:    "https://tls.peet.ws/api/all",
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
				Body:   []byte(tt.args.body),
			})

			if err != nil {
				panic(err)
			}

			fmt.Println(resp.Body)

			resp2, err := c.DoTls(cleanhttp.RequestOption{
				Method: tt.args.method,
				Url:    tt.args.url,
				Header: c.GetDefaultHeader(),
				Body:   []byte(tt.args.body),
			})

			if err != nil {
				panic(err)
			}

			fmt.Println(resp2.Body)

			fmt.Println("=======================================")
		})
	}
}
