package rtc

import (
	"strings"
)

// ClientInfo parses client SDK version, platform, and browser agent string.
type ClientInfo struct {
	SDK        string `json:"sdk"`
	Version    string `json:"version"`
	OS         string `json:"os"`
	Browser    string `json:"browser"`
	Protocol   int    `json:"protocol"`
}

// ParseClientInfo extracts platform details from standard User-Agent header strings.
func ParseClientInfo(userAgent string) ClientInfo {
	info := ClientInfo{
		SDK:      "unknown",
		Protocol: 12,
	}

	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "js") || strings.Contains(ua, "chrome") || strings.Contains(ua, "firefox") {
		info.SDK = "js"
		info.Browser = "browser"
	} else if strings.Contains(ua, "swift") || strings.Contains(ua, "ios") {
		info.SDK = "swift"
		info.OS = "ios"
	} else if strings.Contains(ua, "android") {
		info.SDK = "android"
		info.OS = "android"
	}
	return info
}
