package discordgo

const (
	droidCapabilities      = 125
	droidOS                = "Windows"
	droidOSVersion         = "10"
	droidBrowser           = "Chrome"
	droidBrowserVersion    = "92.0.4515.159"
	droidBrowserUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + droidBrowserVersion + " Safari/537.36"
	droidReferrer          = "https://discord.com/channels/@me"
	droidReferringDomain   = "discord.com"
	droidClientBuildNumber = "83364"
	droidReleaseChannel    = "stable"
	droidStatus            = "online"
)

var (
	droidWSHeaders = map[string]string{
		"User-Agent":    droidBrowserUserAgent,
		"Origin":        "https://discord.com",
		"Pragma":        "no-cache",
		"Cache-Control": "no-cache",
	}
)
