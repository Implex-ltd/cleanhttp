package cleanhttp

import (
	"github.com/bogdanfinn/fhttp/http2"
	tls_client "github.com/bogdanfinn/tls-client"
	tls "github.com/bogdanfinn/utls"
)

// Chrome 114 exact browser TLS.
func GetTlsProfile() tls_client.ClientProfile {
	settings := map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:      65536,
		http2.SettingMaxConcurrentStreams: 1000,
		http2.SettingInitialWindowSize:    6291456,
		http2.SettingMaxHeaderListSize:    262144,
		http2.SettingEnablePush:           0,
	}
	settingsOrder := []http2.SettingID{
		http2.SettingHeaderTableSize,
		http2.SettingEnablePush,
		http2.SettingMaxConcurrentStreams,
		http2.SettingInitialWindowSize,
		http2.SettingMaxHeaderListSize,
	}

	pseudoHeaderOrder := []string{
		":method",
		":authority",
		":scheme",
		":path",
	}

	connectionFlow := uint32(15663105)

	specFactory := func() (tls.ClientHelloSpec, error) {
		return tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				tls.CompressionNone,
			},
			Extensions: []tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.SNIExtension{},
				&tls.UtlsExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(tls.GREASE_PLACEHOLDER),
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
					//25497,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
				}},
				&tls.SCTExtension{},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.GREASE_PLACEHOLDER,
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsCompressCertExtension{Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
				&tls.ApplicationSettingsExtension{SupportedProtocols: []string{"h2"}},
				&tls.UtlsGREASEExtension{},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}, nil
	}

	customClientProfile := tls_client.NewClientProfile(tls.ClientHelloID{
		Client:      "Chrome114",
		Version:     "1",
		Seed:        nil,
		SpecFactory: specFactory,
	}, settings, settingsOrder, pseudoHeaderOrder, connectionFlow, nil, nil)

	return customClientProfile
}
