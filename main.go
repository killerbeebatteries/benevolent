package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST            = "irc.libera.chat"
	CONN_PORT            = "6697"
	CONN_TYPE            = "tcp"
	BOT_NAME             = "benbot"
	SECURE               = true
	USE_CHANNEL_PASSWORD = true
	USE_NICKSERV         = true
)

// IRCBot represents an IRC bot.
type IRCBot struct {
	conn net.Conn
}

// NewIRCBot creates a new instance of IRCBot.
func NewIRCBot(server, port, nickname string, secure bool) (*IRCBot, error) {
	var bot IRCBot
	var err error

	if secure {
		config := &tls.Config{InsecureSkipVerify: true}
		bot.conn, err = tls.Dial(CONN_TYPE, fmt.Sprintf("%s:%s", server, port), config)
	} else {
		bot.conn, err = net.Dial(CONN_TYPE, fmt.Sprintf("%s:%s", server, port))
	}

	if err != nil {
		return nil, err
	}

	// Perform IRC handshake
	err = bot.sendRaw(fmt.Sprintf("NICK %s", nickname))
	if err != nil {
		fmt.Println("Error sending NICK:", err)
	}

	err = bot.sendRaw(fmt.Sprintf("USER %s 0 * :%s", nickname, nickname))
	if err != nil {
		fmt.Println("Error sending USER:", err)
	}

	return &bot, nil
}

// sendRaw sends a raw IRC command to the server.
func (b *IRCBot) sendRaw(command string) error {
	if b.conn == nil {
		return errors.New("connection is nil")
	}

	_, err := fmt.Fprintf(b.conn, "%s\r\n", command)
	if err != nil {
		return err
	}

	return nil
}

// joinChannel joins a specified IRC channel.
func (b *IRCBot) joinChannel(channel string, password string) {
	if password != "" {
		b.sendRaw(fmt.Sprintf("JOIN %s %s", channel, password))
	} else {
		b.sendRaw(fmt.Sprintf("JOIN %s", channel))
	}
}

// sendMessage sends a message to a specified IRC channel.
func (b *IRCBot) sendMessage(channel, message string) {
	b.sendRaw(fmt.Sprintf("PRIVMSG %s :%s", channel, message))
	// throttle messages to avoid being kicked
	time.Sleep(200 * time.Millisecond)
}

// receiveMessages continuously reads and processes messages from the IRC server.
func (b *IRCBot) receiveMessages() {

	// TODO: Need to learn how to declare globals.
	// I don't want to pass it as an argument, as it feels like that is misleading.
	// It leads you to think we're receiving messages from a channel, rather than a server.
	CHANNEL := os.Getenv("CHANNEL")

	if CHANNEL == "" {
		log.Fatal("CHANNEL environment variable not set")
	}

	TRUSTED_USERS := strings.Split(os.Getenv("TRUSTED_USERS"), ",")

	if len(TRUSTED_USERS) == 0 {
		log.Fatal("TRUSTED_USERS environment variable not set")
	}

	scanner := bufio.NewScanner(b.conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received:", message)

		event := strings.Split(message, " ")[1]
		user := strings.ToLower(getUserFromMessage(message))
		userIsTrusted := false

		for _, line := range TRUSTED_USERS {
			if strings.EqualFold(user, line) {
				userIsTrusted = true
			}
		}

		// Add your message processing logic here
		// Example: check for PING messages and respond with PONG
		if strings.HasPrefix(message, "PING") {
			fmt.Println("Sending: PONG " + message[5:])
			b.sendRaw("PONG " + message[5:])
		}

		if event == "JOIN" {

			if user != BOT_NAME {
				resp, err := getGreetings(user)

				if err != nil {
					fmt.Println("Error retrieving greeting message:", err)
				}

				b.sendMessage(CHANNEL, resp)
			}

			if user != BOT_NAME {
				resp, err := sendRelayMessage(user)

				if err != nil {
					fmt.Println("Error sending relay message:", err)
				}

				for _, line := range resp {
					b.sendMessage(CHANNEL, line)
				}
			}
		}

		// case statement for commands
		if userIsTrusted {
			if event == "PRIVMSG" {
				command := strings.Split(message, " ")[3]
				switch command {
				case ":!hello":
					b.sendMessage(CHANNEL, "Hello, world!")
				case ":!ping":
					b.sendMessage(CHANNEL, "pong")
				case ":!time":
					b.sendMessage(CHANNEL, time.Now().String())
				case ":!weather":
					if len(strings.Split(message, " ")) > 4 {
						location := strings.Join(strings.Split(message, " ")[4:], " ")
						fmt.Println("Checking weather for location:", location)

						if forecast, err := handleWeather(location); err != nil {
							fmt.Println("Error getting weather:", err)
						} else {
							fmt.Println("Sending weather forecast:", forecast)
							for _, line := range forecast {
								b.sendMessage(CHANNEL, line)
							}
						}
					} else {
						resp, err := getHelp("weather")
						if err != nil {
							fmt.Println("Error retrieving help for weather: ", err)
						}
						for _, line := range resp {
							b.sendMessage(CHANNEL, line)
						}
					}
				case ":!relay_url":
					if len(strings.Split(message, " ")) > 4 {
						resp, err := addRelayMessage(message)
						if err != nil {
							fmt.Println("Error adding relay message:", err)
						}
						for _, line := range resp {
							b.sendMessage(CHANNEL, line)
						}
					} else {
						resp, err := getHelp("relay_url")
						if err != nil {
							fmt.Println("Error retrieving help for relay url messages:", err)
						}
						for _, line := range resp {
							b.sendMessage(CHANNEL, line)
						}
					}
				case ":!help":
					if len(strings.Split(message, " ")) > 4 {
            feature := strings.Split(message, " ")[4]
            resp, err := getHelp(feature)
						if err != nil {
							fmt.Println("Error retrieving help message:", err)
						}
						for _, line := range resp {
							b.sendMessage(CHANNEL, line)
						}
					} else {
						resp, err := getHelp("")
						if err != nil {
							fmt.Println("Error retrieving general help message:", err)
						}
						for _, line := range resp {
							b.sendMessage(CHANNEL, line)
						}
					}
					// case ":!quit":
					//   b.sendMessage(CHANNEL, "Bye!")
					//   b.sendRaw("QUIT")
					//   b.conn.Close()
					//   return
				}
			}
		} else if event == "PRIVMSG" {
			fmt.Println("User not trusted:", user)
		}
	}
}

func main() {

	err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDatabase()

	CHANNEL := os.Getenv("CHANNEL")
	CHANNEL_PASSWORD := ""

	if CHANNEL == "" {
		log.Fatal("CHANNEL environment variable not set")
	}

	if USE_CHANNEL_PASSWORD {
		CHANNEL_PASSWORD = os.Getenv("CHANNEL_PASSWORD")
		if CHANNEL_PASSWORD == "" {
			log.Fatal("CHANNEL_PASSWORD environment variable not set")
		}
	}

	bot, err := NewIRCBot(CONN_HOST, CONN_PORT, BOT_NAME, SECURE)
	if err != nil {
		fmt.Println("Error creating IRC bot:", err)
		return
	}
	defer bot.conn.Close()

	// Start a goroutine to handle incoming messages
	go bot.receiveMessages()
	// wait for 10 seconds before joining the channel
	<-time.After(10 * time.Second)
	bot.joinChannel(CHANNEL, CHANNEL_PASSWORD)

	if USE_NICKSERV {
		NICKSERV_PASSWORD := os.Getenv("NICKSERV_PASSWORD")
		if NICKSERV_PASSWORD == "" {
			log.Fatal("NICKSERV_PASSWORD environment variable not set")
		}
		bot.sendRaw(fmt.Sprintf("PRIVMSG nickserv :identify %s", NICKSERV_PASSWORD))
	}

	// keep channel alive.
	// TODO: This is likely a bit of a hack. There's probably a better way to do this.
	for {
		<-time.After(10 * time.Second)
	}
	// Example: Send a message to the channel every 10 seconds
	// for {
	// 	bot.sendMessage(CHANNEL, "Hello, IRC!")
	// 	<-time.After(10 * time.Second)
	// }
}
