module gitlab.com/beeper/discordgo

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/gorilla/websocket v1.4.0
	golang.org/x/crypto v0.0.0-20181030102418-4d3f4d9ffa16
)

replace (
	github.com/bwmarrin/discordgo v0.23.2 => github.com/grimmy/discordgo v0.23.3-0.20220126043435-7470d1aacd64
	github.com/bwmarrin/discordgo v0.32.2 => gitlab.com/beeper/discordgo v0.23.3-0.20220127181915-5589d3741f1b
)

go 1.10
