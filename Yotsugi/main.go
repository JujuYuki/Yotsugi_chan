package main

import (
	"sync"
	"time"
	"twitchBot"
)

func main() {
	wg := new(sync.WaitGroup)
	yotsugi := twitchBot.BasicBot{
		Channel:		"megameloetta",
		Owner:			"jujuyuki",
		MsgRate:		time.Duration(20/30)*time.Millisecond,
		Name:			"Yotsugi_chan",
		Port:			"6667",
		Private:		"./private/oauth.json",
		Server:			"irc.chat.twitch.tv",
		WaitGroup:		wg,
	}
	wg.Add(1)
	go yotsugi.Start()
	yotsugi2 := twitchBot.BasicBot{
		Channel:		"jujuyuki",
		Owner:			"jujuyuki",
		MsgRate:		time.Duration(20/30)*time.Millisecond,
		Name:			"Yotsugi_chan",
		Port:			"6667",
		Private:		"./private/oauth.json",
		Server:			"irc.chat.twitch.tv",
		WaitGroup:		wg,
	}
	wg.Add(1)
	go yotsugi2.Start()
	wg.Wait()
}