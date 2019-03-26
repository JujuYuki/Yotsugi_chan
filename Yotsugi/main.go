package main

import (
	"time"
	"twitchBot"
)

func main() {
	yotsugi := twitchBot.BasicBot{
		Channel:		"megameloetta",
		Owner:			"jujuyuki",
		MsgRate:		time.Duration(20/30)*time.Millisecond,
		Name:			"Yotsugi_chan",
		Port:			"6667",
		Private:		"./private/oauth.json",
		Server:			"irc.chat.twitch.tv",
	}
	yotsugi.Start()
	yotsugi2 := twitchBot.BasicBot{
		Channel:		"jujuyuki",
		Owner:			"jujuyuki",
		MsgRate:		time.Duration(20/30)*time.Millisecond,
		Name:			"Yotsugi_chan",
		Port:			"6667",
		Private:		"./private/oauth.json",
		Server:			"irc.chat.twitch.tv",
	}
	yotsugi2.Start()
}