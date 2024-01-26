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
func NewIRCBot(server, port, nickname, channel string) (*IRCBot, error) {
  // Establish a TCP connection to the IRC server
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", server, port))
  if err != nil {
    return nil, err
  }

  bot := &IRCBot{conn: conn}

  // Perform IRC handshake and join the channel
  bot.sendRaw(fmt.Sprintf("NICK %s", nickname))
  bot.sendRaw(fmt.Sprintf("USER %s 0 * :%s", nickname, nickname))
  bot.joinChannel(channel)

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
      b.sendRaw("PONG " + message[5:])
    }
  }
}

func main() {
  // Replace these with your IRC server details
  server := CONN_HOST
  port := CONN_PORT
  nickname := BOT_NAME
  channel := CHANNEL

  bot, err := NewIRCBot(server, port, nickname, channel)
  if err != nil {
    fmt.Println("Error creating IRC bot:", err)
    return
  }
  defer bot.conn.Close()

  go bot.receiveMessages() // Start a goroutine to handle incoming messages

  // Example: Send a message to the channel every 10 seconds
  for {
    bot.sendMessage(channel, "Hello, IRC!")
    <-time.After(10 * time.Second)
  }
}
