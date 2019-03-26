package twitchBot

import (
	"fmt"
	"net"
	"time"
)

const UTCFormat = "Jan 2 15:04:05 UTC"

type BasicBot struct {
	Channel		string
	conn		net.Conn
	Credentials	*OAuthCred
	MsgRate		time.Duration
	Name		string
	Port		string
	Private		string
	Server		string
	startTime	time.Time
}

type OAuthCred struct {
	Password 	string `json:"password,omitempty"`
}

type TwitchBot interface {
	Connect()
	Disconnect()
	HandleChat() error
	JoinChannel()
	ReadCredentials() (*OAuthCred, error)
	Say(msg string) error
	Start()
}

func timeStamp() string {
	return TimeStamp(UTCFormat)
}

func TimeStamp(fmt string) string {
	return time.Now().Format(fmt)
}

func (bot *BasicBot) Connect() {
	var err error
	fmt.Printf("[%s] Connecting: %s...\n", timeStamp(), bot.Server)

	bot.conn, err = net.Dial("tcp", bot.Server+":"+bot.Port)
	if nil != err {
		fmt.Printf("[%s] Cannot connect to %s, retrying.\n", timeStamp(), bot.Server)
		bot.Connect()
		return
	}
	fmt.Printf("[%s] Connected: %s!\n", timeStamp(), bot.Server)
	bot.startTime = time.Now()
}

func (bot *BasicBot) Disconnect() {
	bot.conn.Close()
	upTime := time.Now().Sub(bot.startTime).Seconds()
	fmt.Printf("[%s] Closed connection to %s, live for %fs\n", timeStamp(), bot.Server, upTime)
}