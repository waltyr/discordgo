package discordgo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	droidCapabilities      = 16381
	droidOS                = "Windows"
	droidOSVersion         = "10"
	droidBrowser           = "Chrome"
	droidReferrer          = "https://discord.com/channels/@me"
	droidReferringDomain   = "discord.com"
	droidClientBuildNumber = "236850"
	droidReleaseChannel    = "stable"
	droidStatus            = "invisible"
	droidSystemLocale      = "en-US"
)

const (
	DroidBrowserMajorVersion = "118"
	DroidBrowserVersion      = DroidBrowserMajorVersion + ".0.0.0"
	DroidBrowserUserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + DroidBrowserVersion + " Safari/537.36"
)

type UserIdentifyProperties struct {
	OS                     string  `json:"os"`
	Browser                string  `json:"browser"`
	Device                 string  `json:"device"`
	SystemLocale           string  `json:"system_locale"`
	BrowserUserAgent       string  `json:"browser_user_agent"`
	BrowserVersion         string  `json:"browser_version"`
	OSVersion              string  `json:"os_version"`
	Referrer               string  `json:"referrer"`
	ReferringDomain        string  `json:"referring_domain"`
	ReferrerCurrent        string  `json:"referrer_current"`
	ReferringDomainCurrent string  `json:"referring_domain_current"`
	ReleaseChannel         string  `json:"release_channel"`
	ClientBuildNumber      string  `json:"client_build_number"`
	ClientEventSource      *string `json:"client_event_source"`
	DesignID               int     `json:"design_id"`
}

type ClientState struct {
	GuildVersions            struct{} `json:"guild_versions"`
	HighestLastMessageID     string   `json:"highest_last_message_id"`
	ReadStateVersion         int      `json:"read_state_version"`
	UserGuildSettingsVersion int      `json:"user_guild_settings_version"`
	UserSettingsVersion      int      `json:"user_settings_version"`
	PrivateChannelsVersion   string   `json:"private_channels_version"`
	APICodeVersion           int      `json:"api_code_version"`
}

func mustMarshalJSON(data interface{}) string {
	dat, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(dat)
}

func basedOn(base map[string]string, additional map[string]string) map[string]string {
	for k, v := range base {
		_, exists := additional[k]
		if !exists {
			additional[k] = v
		}
	}
	return additional
}

var (
	droidIdentifyProperties = &UserIdentifyProperties{
		OS:               droidOS,
		OSVersion:        droidOSVersion,
		Browser:          droidBrowser,
		BrowserVersion:   DroidBrowserVersion,
		BrowserUserAgent: DroidBrowserUserAgent,
		//Referrer: droidReferrer,
		//ReferringDomain: droidReferringDomain,
		ClientBuildNumber: droidClientBuildNumber,
		ReleaseChannel:    droidReleaseChannel,
		SystemLocale:      droidSystemLocale,
		DesignID:          0,
	}
	DroidFetchHeaders = map[string]string{
		"Sec-CH-UA":          fmt.Sprintf(`" Not A;Brand";v="99", "Chromium";v="%[1]s", "Google Chrome";v="%[1]s"`, DroidBrowserMajorVersion),
		"Sec-CH-UA-Mobile":   "?0",
		"Sec-CH-UA-Platform": droidOS,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "sameorigin",
		"X-Discord-Locale":   droidSystemLocale,
		"X-Super-Properties": mustMarshalJSON(droidIdentifyProperties),

		"Accept":          "*/*",
		"Origin":          "https://discord.com",
		"Accept-Language": "en-US,en;q=0.9",
		"User-Agent":      DroidBrowserUserAgent,
	}
	DroidDownloadHeaders = basedOn(DroidFetchHeaders, map[string]string{
		"Sec-Fetch-Mode": "no-cors",
	})
	DroidImageHeaders = basedOn(DroidDownloadHeaders, map[string]string{
		"Accept":         "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
		"Sec-Fetch-Dest": "image",
	})

	DroidWSHeaders = map[string]string{
		"User-Agent":      DroidBrowserUserAgent,
		"Origin":          "https://discord.com",
		"Accept-Language": "en-US,en;q=0.9",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Accept-Encoding": "gzip, deflate, br",

		//"Sec-Fetch-Dest": "websocket",
		//"Sec-Fetch-Mode": "websocket",
		//"Sec-Fetch-Site": "cross-site",
	}
)

const (
	ThreadJoinLocationContextMenu     = "Context Menu"
	ThreadJoinLocationToolbarOverflow = "Toolbar Overflow"
	ThreadJoinLocationSidebarOverflow = "Sidebar Overflow"
)

func (s *Session) ThreadJoinWithLocation(id, location string) error {
	if !s.IsUser {
		return s.ThreadJoin(id)
	}
	endpoint := EndpointThreadMember(id, "@me") + "?" + url.Values{
		"location": []string{location},
	}.Encode()
	_, err := s.RequestWithBucketID("PUT", endpoint, nil, endpoint)
	return err
}
