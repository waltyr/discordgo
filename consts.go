package discordgo

const (
	droidCapabilities      = 509
	droidOS                = "Windows"
	droidOSVersion         = "10"
	droidBrowser           = "Chrome"
	droidReferrer          = "https://discord.com/channels/@me"
	droidReferringDomain   = "discord.com"
	droidClientBuildNumber = "130153"
	droidReleaseChannel    = "stable"
	droidStatus            = "online"
	droidSystemLocale      = "en-US"
)

const (
	DroidBrowserVersion   = "102.0.5005.61"
	DroidBrowserUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + DroidBrowserVersion + " Safari/537.36"
)

var (
	droidWSHeaders = map[string]string{
		"User-Agent":    DroidBrowserUserAgent,
		"Origin":        "https://discord.com",
		"Pragma":        "no-cache",
		"Cache-Control": "no-cache",

		//"Sec-Fetch-Dest": "websocket",
		//"Sec-Fetch-Mode": "websocket",
		//"Sec-Fetch-Site": "cross-site",
	}
)
