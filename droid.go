package discordgo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	droidOS              = "Windows"
	droidOSVersion       = "10"
	droidBrowser         = "Chrome"
	droidReferrer        = "https://discord.com/channels/@me"
	droidReferringDomain = "discord.com"
	droidReleaseChannel  = "stable"
	droidStatus          = "invisible"
	droidSystemLocale    = "en-US"
)

var (
	droidCapabilities      = 30717
	droidClientBuildNumber = 370533
	droidGatewayURL        = ""
	mainPageLoaded         = false
)

var mainPageLoadLock sync.Mutex

const (
	DroidBrowserMajorVersion = "133"
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
	ClientBuildNumber      int     `json:"client_build_number"`
	ClientEventSource      *string `json:"client_event_source"`
}

type ClientState struct {
	GuildVersions            struct{} `json:"guild_versions"`
	HighestLastMessageID     string   `json:"highest_last_message_id,omitempty"`
	ReadStateVersion         int      `json:"read_state_version,omitempty"`
	UserGuildSettingsVersion int      `json:"user_guild_settings_version,omitempty"`
	UserSettingsVersion      int      `json:"user_settings_version,omitempty"`
	PrivateChannelsVersion   string   `json:"private_channels_version,omitempty"`
	APICodeVersion           int      `json:"api_code_version,omitempty"`
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

func UpdateVersion(version, capabilities int) {
	droidClientBuildNumber = version
	droidCapabilities = capabilities
	droidIdentifyProperties.ClientBuildNumber = version
	DroidFetchHeaders["X-Super-Properties"] = mustMarshalJSON(droidIdentifyProperties)
	DroidDownloadHeaders["X-Super-Properties"] = DroidFetchHeaders["X-Super-Properties"]
	DroidImageHeaders["X-Super-Properties"] = DroidFetchHeaders["X-Super-Properties"]
}

func (s *Session) SetGatewayURL(url string) {
	s.gateway = url + "?encoding=json&v=" + APIVersion + "&compress=zlib-stream"
	s.noClearGateway = true
}

var apiVersionRegex = regexp.MustCompile(`API_VERSION: (\d+),`)
var gatewayURLRegex = regexp.MustCompile(`GATEWAY_ENDPOINT:\s?['"](.+?)['"],`)
var mainJSRegex = regexp.MustCompile(`src="(/assets/web.[a-f0-9]{12,32}.js)"`)
var buildNumberRegex = regexp.MustCompile(`(?:buildNumber|build_number):\s?['"]?(\d{6,})['"]?`)

func (s *Session) LoadMainPage(ctx context.Context) error {
	mainPageLoadLock.Lock()
	defer mainPageLoadLock.Unlock()
	if mainPageLoaded && droidGatewayURL != "" {
		s.SetGatewayURL(droidGatewayURL)
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://discord.com/channels/@me", nil)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}
	for name, value := range DroidBaseHeaders {
		req.Header.Add(name, value)
	}
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch main page: %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read main page: %w", err)
	}

	apiVersionMatch := apiVersionRegex.FindSubmatch(data)
	if apiVersionMatch == nil {
		return fmt.Errorf("failed to find API version")
	} else if string(apiVersionMatch[1]) != APIVersion {
		return fmt.Errorf("API version mismatch: expected %s, got %s", APIVersion, apiVersionMatch[1])
	}
	gatewayURLMatch := gatewayURLRegex.FindSubmatch(data)
	if gatewayURLMatch == nil {
		return fmt.Errorf("failed to find gateway URL")
	}
	droidGatewayURL = string(gatewayURLMatch[1])
	if !strings.HasSuffix(droidGatewayURL, "/") {
		droidGatewayURL += "/"
	}
	s.log(LogInformational, "Found gateway URL %s and confirmed API version", droidGatewayURL)
	s.SetGatewayURL(droidGatewayURL)
	mainJSMatch := mainJSRegex.FindSubmatch(data)
	if mainJSMatch == nil {
		return fmt.Errorf("failed to find main JS URL")
	}

	jsReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://discord.com"+string(mainJSMatch[1]), nil)
	if err != nil {
		return fmt.Errorf("failed to prepare JS request: %w", err)
	}
	for name, value := range DroidBaseHeaders {
		req.Header.Add(name, value)
	}
	jsReq.Header.Set("Sec-Fetch-Dest", "script")
	jsReq.Header.Set("Sec-Fetch-Mode", "no-cors")
	jsReq.Header.Set("Sec-Fetch-Site", "same-origin")
	jsReq.Header.Set("Accept", "*/*")
	jsResp, err := s.Client.Do(jsReq)
	if err != nil {
		return fmt.Errorf("failed to fetch JS: %w", err)
	}
	jsData, err := io.ReadAll(jsResp.Body)
	_ = jsResp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read JS: %w", err)
	}
	buildNumberMatch := buildNumberRegex.FindSubmatch(jsData)
	if buildNumberMatch == nil {
		return fmt.Errorf("failed to find build number")
	}
	buildNumberInt, err := strconv.Atoi(string(buildNumberMatch[1]))
	if err != nil {
		return fmt.Errorf("failed to parse build number %s: %w", buildNumberMatch[1], err)
	}
	s.log(LogInformational, "Found build number %d from JS file %s", buildNumberInt, string(mainJSMatch[1]))
	// TODO parse capabilities too?
	UpdateVersion(buildNumberInt, droidCapabilities)
	mainPageLoaded = true

	return nil
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
	}
	DroidBaseHeaders = map[string]string{
		"Sec-Ch-Ua":          fmt.Sprintf(`" Not A;Brand";v="99", "Chromium";v="%[1]s", "Google Chrome";v="%[1]s"`, DroidBrowserMajorVersion),
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": `"` + droidOS + `"`,

		"Accept":          "*/*",
		"Origin":          "https://discord.com",
		"Accept-Language": "en-US,en;q=0.9",
		"User-Agent":      DroidBrowserUserAgent,
	}
	DroidFetchHeaders = basedOn(DroidBaseHeaders, map[string]string{
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"X-Debug-Options":    "bugReporterEnabled",
		"X-Discord-Locale":   droidSystemLocale,
		"X-Discord-Timezone": "UTC",
		"X-Super-Properties": mustMarshalJSON(droidIdentifyProperties),
	})
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

const (
	ReactionLocationHoverBar     = "Message Hover Bar"
	ReactionLocationInlineButton = "Message Inline Button"
	ReactionLocationPicker       = "Message Reaction Picker"
	ReactionLocationContextMenu  = "Message Context Menu"
)

func (s *Session) MessageReactionAddUser(guildID, channelID, messageID, emojiID string, options ...RequestOption) error {
	if s.IsUser {
		options = append(
			options,
			WithChannelReferer(guildID, channelID),
			WithLocationParam(ReactionLocationPicker),
			WithQueryParam("type", "0"),
		)
	}
	return s.MessageReactionAdd(channelID, messageID, emojiID, options...)
}

func (s *Session) MessageReactionRemoveUser(guildID, channelID, messageID, emojiID, userID string, options ...RequestOption) error {
	if s.IsUser {
		options = append(
			options,
			WithChannelReferer(guildID, channelID),
			WithLocationParam(ReactionLocationInlineButton),
			WithQueryParam("burst", "false"),
		)
	}
	return s.MessageReactionRemove(channelID, messageID, emojiID, userID, options...)
}
