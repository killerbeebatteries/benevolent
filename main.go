package main

import (
  "bufio"
  "fmt"
  "net"
  "strings"
  "time"
)

const (
  CONN_HOST = "localhost"
  CONN_PORT = "6667"
  CONN_TYPE = "tcp"
  BOT_NAME  = "benbot"
  CHANNEL   = "#lurking"
)

// IRCBot represents the IRC bot structure.
type IRCBot struct {
  conn net.Conn
}

// NewIRCBot creates a new instance of IRCBot.
func NewIRCBot(server, port, nickname string) (*IRCBot, error) {
  // Establish a TCP connection to the IRC server
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server, port))
  if err != nil {
    return nil, err
  }

  bot := &IRCBot{conn: conn}

  // Perform IRC handshake and join the channel
  bot.sendRaw(fmt.Sprintf("NICK %s", nickname))
  bot.sendRaw(fmt.Sprintf("USER %s 0 * :%s", nickname, nickname))

  return bot, nil
}

// sendRaw sends a raw IRC command to the server.
func (b *IRCBot) sendRaw(command string) {
  fmt.Fprintf(b.conn, "%s\r\n", command)
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

    // Add your message processing logic here
    // Example: check for PING messages and respond with PONG
    if strings.HasPrefix(message, "PING") {
      fmt.Println("Sending: PONG " + message[5:])
      b.sendRaw("PONG " + message[5:])
    }

    // case statement for commands
    if strings.Contains(message, "PRIVMSG") {
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
          b.sendMessage(CHANNEL, "Usage: !weather <location>")
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

  bot, err := NewIRCBot(CONN_HOST, CONN_PORT, BOT_NAME)
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

  // Example: Send a message to the channel every 10 seconds
  for {
   bot.sendMessage(CHANNEL, "Hello, IRC!")
   <-time.After(10 * time.Second)
  }
}
