package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"github.com/Implex-ltd/fingerprint-client/fpclient"
	http "github.com/bogdanfinn/fhttp"
)

func main() {
	fp, err := fpclient.LoadFingerprint(&fpclient.LoadingConfig{
		FilePath: "../assets/chrome114.json",
	})

	if err != nil {
		panic(err)
	}

	c, err := cleanhttp.NewCleanHttpClient(&cleanhttp.Config{
		BrowserFp: fp,
		//Proxy:     "http://user:pass@ip:port",
	})

	if err != nil {
		panic(err)
	}

	/**
	 * Return these params based on your fingerprint:
	 * 	- cookies
	 * 	- SecChUa
	 * 	- secChUaMobile
	 * 	- AcceptLanguage
	 * 	- SecChUaPlatform
	 */
	base := c.GenerateBaseHeaders()

	/**
	 * Into this example, the header is custom one.
	 * To get default one use: c.GetDefaultHeader()
	 */

	resp, err := c.Do(cleanhttp.RequestOption{
		Method: "GET", // GET, POST, PUT, PATCH, DELETE
		Url:    "https://discord.com/api/v9/experiments",
		Header: http.Header{
			`accept`:             {`*/*`},
			`accept-encoding`:    {`gzip, deflate, br`},
			`accept-language`:    {base.AcceptLanguage},
			`content-type`:       {`application/json`},
			`cookie`:             {base.Cookies},
			`origin`:             {`https://google.com`},
			`referer`:            {`https://google.com`},
			`sec-ch-ua`:          {base.SecChUa},
			`sec-ch-ua-mobile`:   {`?0`},
			`sec-ch-ua-platform`: {base.SecChUaPlatform},
			`sec-fetch-dest`:     {`empty`},
			`sec-fetch-mode`:     {`cors`},
			`sec-fetch-site`:     {`same-origin`},
			`user-agent`:         {c.Config.BrowserFp.Navigator.UserAgent},

			// Keep your headers in the right order.
			http.HeaderOrderKey: {
				`authority`,
				`accept`,
				`accept-encoding`,
				`accept-language`,
				`content-type`,
				`cookie`,
				`origin`,
				`referer`,
				`sec-ch-ua`,
				`sec-ch-ua-mobile`,
				`sec-ch-ua-platform`,
				`sec-fetch-dest`,
				`sec-fetch-mode`,
				`sec-fetch-site`,
				`user-agent`,
			},
		},
		//Body: []byte(`{"content": "this is my super json string !"}`),
	})

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Exact same TLS as chrome 114 !
	resp, err = c.Do(cleanhttp.RequestOption{
		Method: "POST", // GET, POST, PUT, PATCH, DELETE
		Url:    "https://discord.com/api/v9/track/ott",
		Header: http.Header{
			`accept`:             {`*/*`},
			`accept-encoding`:    {`gzip, deflate, br`},
			`accept-language`:    {base.AcceptLanguage},
			`content-type`:       {`application/json`},
			`cookie`:             {base.Cookies},
			`origin`:             {`https://google.com`},
			`referer`:            {`https://google.com`},
			`sec-ch-ua`:          {base.SecChUa},
			`sec-ch-ua-mobile`:   {`?0`},
			`sec-ch-ua-platform`: {base.SecChUaPlatform},
			`sec-fetch-dest`:     {`empty`},
			`sec-fetch-mode`:     {`cors`},
			`sec-fetch-site`:     {`same-origin`},
			`user-agent`:         {c.Config.BrowserFp.Navigator.UserAgent},

			// Keep your headers in the right order.
			http.HeaderOrderKey: {
				`authority`,
				`accept`,
				`accept-encoding`,
				`accept-language`,
				`content-type`,
				`cookie`,
				`origin`,
				`referer`,
				`sec-ch-ua`,
				`sec-ch-ua-mobile`,
				`sec-ch-ua-platform`,
				`sec-fetch-dest`,
				`sec-fetch-mode`,
				`sec-fetch-site`,
				`user-agent`,
			},
		},
		Body: strings.NewReader(`{"type":"landing"}`),
	})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// cookies have been set.
	fmt.Println(resp.Request.Header)
	fmt.Println(string(data))
}
