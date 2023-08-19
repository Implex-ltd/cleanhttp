package cleanhttp

import (
	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"reflect"
	"testing"
)

func TestParseUserAgent(t *testing.T) {
	type args struct {
		userAgentString string
	}
	tests := []struct {
		name string
		args args
		want *cleanhttp.UserAgentInfo
	}{
		{
			name: "chrome114",
			args: args{
				userAgentString: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
			},
			want: &cleanhttp.UserAgentInfo{
				BrowserName:    "Chrome",
				BrowserVersion: "114.0.0.0",
				OSName:         "Windows",
				OSVersion:      "10",
				UaVersion:      "114",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanhttp.ParseUserAgent(tt.args.userAgentString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUserAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
