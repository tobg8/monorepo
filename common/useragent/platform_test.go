package useragent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlatformFromUserAgent(t *testing.T) {
	t.Run("with user agent", func(t *testing.T) {
		for userAgent, expectedPlateform := range map[string]Platform{
			"Adevinta;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                      {Kind: PlatformKindAndroidWifi, Version: "4.45.9.0"},
			"Kleinanzeigen;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                 {Kind: PlatformKindAndroidWifi, Version: "4.45.9.0"},
			"LBC;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                           {Kind: PlatformKindAndroidWifi, Version: "4.45.9.0"},
			"LBC;Android;9;SM-A105FN;phone;806db1fdd1b36de0;wwan;4.44.3.0;84300;0":                                                                                           {Kind: PlatformKindAndroidWWAN, Version: "4.44.3.0"},
			"LBC;Android;8.0.0;BAH2-L09;tab;05c9b52af1589e17;wifi;4.44.3.0;84300;0":                                                                                          {Kind: PlatformKindAndroidTabWifi, Version: "4.44.3.0"},
			"LBC;Android;5.1.1;SM-T285;tab;ced308ba41534c5b;wwan;4.40.3.0;80300;0":                                                                                           {Kind: PlatformKindAndroidTabWWAN, Version: "4.40.3.0"},
			"LBC;Android;5.0;E2303;phone;19e04ff3bf99e2d4;;4.7.2.0;47200;-1":                                                                                                 {Kind: PlatformKindAndroidUnknown, Version: "4.7.2.0"},
			"LBC;iOS;13.2.3;iPhone;phone;656090CC-DCEE-484A-A289-6689AA8FDF8C;wifi;5.1.0;201912141044.624":                                                                   {Kind: PlatformKindIOSWifi, Version: "5.1.0"},
			"LBC;iOS;13.3;iPhone;phone;275A0186-5BFE-4D61-8764-6522CBFED59F;wwan;5.1.0;201912141044.624":                                                                     {Kind: PlatformKindIOSWWAN, Version: "5.1.0"},
			"LBC;iOS;12.2;iPad;tab;62E3646F-A7B8-4295-970A-9668E3AC6BCD;wifi;5.1.0;201912141044.624":                                                                         {Kind: PlatformKindIOSTabWifi, Version: "5.1.0"},
			"LBC;iOS;13.3;iPad;tab;5AE4C85A-DDF6-4410-B436-E9D20A797A5B;wwan;5.1.0;201912141044.624":                                                                         {Kind: PlatformKindIOSTabWWAN, Version: "5.1.0"},
			"Mozilla/5.0 (Linux; Android 9; SM-G965F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.93 Mobile Safari/537.36": {Kind: PlatformKindResponsive, Version: ""},
			"PostmanRuntime/7.29.0": {Kind: PlatformKindResponsive, Version: ""},
			"curl/7.54.1":           {Kind: PlatformKindResponsive, Version: ""},
		} {
			platform := PlatformFromUserAgent(userAgent)
			assert.Equal(t, expectedPlateform, platform)
		}
	})

	t.Run("without user agent", func(t *testing.T) {
		platform := PlatformFromUserAgent("")
		assert.Empty(t, platform)
	})
}

func Test_getPlatformVersion(t *testing.T) {
	for userAgent, expectedVersion := range map[string]string{
		"Adevinta;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                      "4.45.9.0",
		"Kleinanzeigen;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                 "4.45.9.0",
		"LBC;Android;10;SM-G975F;phone;3ed3fb71a72871e4;wifi;4.45.9.0;85900;0":                                                                                           "4.45.9.0",
		"LBC;Android;9;SM-A105FN;phone;806db1fdd1b36de0;wwan;4.44.3.0;84300;0":                                                                                           "4.44.3.0",
		"LBC;Android;8.0.0;BAH2-L09;tab;05c9b52af1589e17;wifi;4.44.3.0;84300;0":                                                                                          "4.44.3.0",
		"LBC;Android;5.1.1;SM-T285;tab;ced308ba41534c5b;wwan;4.40.3.0;80300;0":                                                                                           "4.40.3.0",
		"LBC;Android;5.0;E2303;phone;19e04ff3bf99e2d4;;4.7.2.0;47200;-1":                                                                                                 "4.7.2.0",
		"LBC;iOS;13.2.3;iPhone;phone;656090CC-DCEE-484A-A289-6689AA8FDF8C;wifi;5.1.0;201912141044.624":                                                                   "5.1.0",
		"LBC;iOS;13.3;iPhone;phone;275A0186-5BFE-4D61-8764-6522CBFED59F;wwan;5.1.0;201912141044.624":                                                                     "5.1.0",
		"LBC;iOS;12.2;iPad;tab;62E3646F-A7B8-4295-970A-9668E3AC6BCD;wifi;5.1.0;201912141044.624":                                                                         "5.1.0",
		"LBC;iOS;13.3;iPad;tab;5AE4C85A-DDF6-4410-B436-E9D20A797A5B;wwan;5.1.0;201912141044.624":                                                                         "5.1.0",
		"Mozilla/5.0 (Linux; Android 9; SM-G965F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.93 Mobile Safari/537.36": "",
		"PostmanRuntime/7.29.0": "",
		"curl/7.54.1":           "",
	} {
		version := getPlatformVersion(userAgent)
		assert.Equal(t, expectedVersion, version)
	}
}
