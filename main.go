package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6697"
	CONN_TYPE = "tcp"
	BOT_NAME  = "benbot"
	CHANNEL   = "#lurking"
	SECURE    = true
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

	// Perform IRC handshake and join the channel
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
func (b *IRCBot) joinChannel(channel string) {
	b.sendRaw(fmt.Sprintf("JOIN %s", channel))
}

// sendMessage sends a message to a specified IRC channel.
func (b *IRCBot) sendMessage(channel, message string) {
	b.sendRaw(fmt.Sprintf("PRIVMSG %s :%s", channel, message))
	// throttle messages to avoid being kicked
	time.Sleep(200 * time.Millisecond)
}

// receiveMessages continuously reads and processes messages from the IRC server.
func (b *IRCBot) receiveMessages() {
	scanner := bufio.NewScanner(b.conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received:", message)

    event := strings.Split(message, " ")[1]

		// Add your message processing logic here
		// Example: check for PING messages and respond with PONG
		if strings.HasPrefix(message, "PING") {
			fmt.Println("Sending: PONG " + message[5:])
			b.sendRaw("PONG " + message[5:])
		}

    if event == "JOIN" {
      user := getUserFromMessage(message)

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
      // case ":!quit":
      //   b.sendMessage(CHANNEL, "Bye!")
      //   b.sendRaw("QUIT")
      //   b.conn.Close()
      //   return
			}
		}
	}
}

func main() {

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
	bot.joinChannel(CHANNEL)

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
