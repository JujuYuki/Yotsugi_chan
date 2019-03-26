package twitchBot

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

// Time format
const UTCFormat = "Jan 2 15:04:05 UTC"

// regex for messages
// messages first
var msgRegex *regexp.Regexp = regexp.MustCompile(`^:(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)
// commands second
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

type BasicBot struct {
	Channel		string
	Owner		string
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
	ReadCredentials() error
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
	fmt.Printf("[%s] Yay, %s! Peace peace!\n", timeStamp(), bot.Server)
	bot.startTime = time.Now()
}

func (bot *BasicBot) Disconnect() {
	bot.conn.Close()
	upTime := time.Now().Sub(bot.startTime).Seconds()
	fmt.Printf("[%s] Closed connection to %s, live for %fs\n", timeStamp(), bot.Server, upTime)
}

func (bot *BasicBot) HandleChat() error {
	fmt.Printf("[%s] Reading #%s...\n", timeStamp(), bot.Channel)

	tp := textproto.NewReader(bufio.NewReader(bot.conn))

	for {
		line, err := tp.ReadLine()
		if nil != err {
			bot.Disconnect()

			return errors.New("bot.HandleChat: Failed to read line from channel. Disconnecting")
		}

		fmt.Printf("[%s] %s\n", timeStamp(), line)

		if "PING :tmi.twitch.tv" == line {
			bot.conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		} else {
			matches := msgRegex.FindStringSubmatch(line)
			if nil != matches {
				userName := matches[1]
				msgType := matches[2]

				switch msgType {
				case "PRIVMSG":
					msg := matches[3]
					fmt.Printf("[%s] %s: %s\n", timeStamp(), userName, msg)

					cmdMatches := cmdRegex.FindStringSubmatch(msg)
					if nil != cmdMatches {
						cmd := cmdMatches[1]

						if userName == bot.Channel || userName == bot.Owner {
							switch cmd {
							case "urboff":
								fmt.Printf("[%s] Unlimited Rulebook: Shutdown.\n", timeStamp())
								bot.Disconnect()
								return nil
							case "test":
								err = bot.Say("Peace peace!")
								if nil != err {
									return err
								}
							default:
								// NOOP
							}
						}
					}

					if strings.Contains(msg, "Yotsugi") {
						err = bot.Say("Yay, onii-chan! Peace peace!")
						if nil != err {
							return err
						}
					}
				default:
					// NOOP
				}
			}
		}
		time.Sleep(bot.MsgRate)
	}
}

func (bot *BasicBot) JoinChannel() {
	fmt.Printf("[%s] Joining #%s...\n", timeStamp(), bot.Channel)
	bot.conn.Write([]byte("PASS "+bot.Credentials.Password+"\r\n"))
	bot.conn.Write([]byte("NICK "+bot.Name+"\r\n"))
	bot.conn.Write([]byte("JOIN #"+bot.Channel+"\r\n"))

	fmt.Printf("[%s] %s has made contact with #%s! Peace peace!", timeStamp(), bot.Name, bot.Channel)
}

func (bot *BasicBot) ReadCredentials() error {
	credFile, err := ioutil.ReadFile(bot.Private)
	if nil != err {
		return err
	}

	bot.Credentials = &OAuthCred{}

	dec := json.NewDecoder(strings.NewReader(string(credFile)))
	if err = dec.Decode(bot.Credentials); nil != err && io.EOF != err {
		return err
	}

	return nil
}

func (bot *BasicBot) Say(msg string) error {
	if "" == msg {
		return errors.New("BasicBot.Say: msg was empty")
	}
	fmt.Println("sending message")
	_, err := bot.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s :%s\r\n", bot.Channel, msg)))
	if nil != err {
		return err
	}
	return nil
}

func (bot *BasicBot) Start() {
	err := bot.ReadCredentials()
	if nil != err {
		fmt.Println(err)
		fmt.Println("Aborting.")
		return
	}

	for  {
		bot.Connect()
		bot.JoinChannel()
		err = bot.HandleChat()
		if nil != err {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err)
			fmt.Printf("Starting %s again.", bot.Name)
		} else {
			return
		}
	}
}