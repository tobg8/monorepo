package useragent

import (
	"regexp"
	"strings"
)

// Platform describes the platform, the http client which execute the http request.
type Platform struct {
	Kind    string
	Version string
}

func (p Platform) String() string { return p.Kind }

// Device returns the platform device kind.
func (p Platform) Device() string {
	switch pf := p.String(); {
	case strings.HasPrefix(pf, "ios_"):
		return "ios"
	case strings.HasPrefix(pf, "android_"):
		return "android"
	default:
		return "browser"
	}
}

// List of supported platform kinds.
const (
	PlatformKindAndroidWifi    = "android_wifi"
	PlatformKindAndroidWWAN    = "android_wwan"
	PlatformKindAndroidTabWifi = "android_tab_wifi"
	PlatformKindAndroidTabWWAN = "android_tab_wwan"
	PlatformKindAndroidUnknown = "android_unknown"
	PlatformKindIOSWifi        = "ios_phone_wifi"
	PlatformKindIOSWWAN        = "ios_phone_wwan"
	PlatformKindIOSTabWifi     = "ios_tab_wifi"
	PlatformKindIOSTabWWAN     = "ios_tab_wwan"
	PlatformKindIOSUnknown     = "ios_unknown"

	// The default value when the platform is not recognized
	PlatformKindResponsive = "responsive"
)

var (
	//PlatformRegex ...
	PlatformRegex = regexp.MustCompile(
		`^(?:LBC|Kleinanzeigen|Adevinta);(iOS|Android);[^;]+;[^;]+;([^;]+);[^;]+;([^;]*)`,
	)

	//TrPlatforms ...
	TrPlatforms = map[string]map[string]map[string]string{
		"iOS": {
			"phone": {
				"wifi": PlatformKindIOSWifi,
				"wwan": PlatformKindIOSWWAN,
			},
			"tab": {
				"wifi": PlatformKindIOSTabWifi,
				"wwan": PlatformKindIOSTabWWAN,
			},
		},
		"Android": {
			"phone": {
				"wifi": PlatformKindAndroidWifi,
				"wwan": PlatformKindAndroidWWAN,
			},
			"tab": {
				"wifi": PlatformKindAndroidTabWifi,
				"wwan": PlatformKindAndroidTabWWAN,
			},
		},
	}
)

func getPlatformVersion(ua string) string {
	parts := strings.Split(ua, ";")
	if len(parts) >= 8 {
		return parts[7]
	}
	return ""
}

// PlatformFromUserAgent returns the platform based on user agent.
func PlatformFromUserAgent(userAgent string) Platform {
	if userAgent != "" {
		matches := PlatformRegex.FindStringSubmatch(userAgent)
		if len(matches) == 4 {
			var (
				appKind    = matches[1]
				mobileKind = matches[2]
				mobileConn = matches[3]
			)
			v, ok := TrPlatforms[appKind][mobileKind][mobileConn]
			if ok {
				return Platform{Kind: v, Version: getPlatformVersion(userAgent)}
			}
			switch appKind {
			case "iOS":
				return Platform{
					Kind:    PlatformKindIOSUnknown,
					Version: getPlatformVersion(userAgent),
				}
			case "Android":
				return Platform{
					Kind:    PlatformKindAndroidUnknown,
					Version: getPlatformVersion(userAgent),
				}
			}
		}
		return Platform{Kind: PlatformKindResponsive}
	}
	return Platform{}
}
