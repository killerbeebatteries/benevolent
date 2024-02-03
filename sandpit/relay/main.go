package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type relayMessage struct {
	Id          int
	timestamp   time.Time
	FromUser    string
	ToUser      string
	Description string
	URL         url.URL
}

func main() {

	var testMessage, fromUser, toUser, channel, description, url, help string
	message = ":alice!unknown@localhost PRIVMSG #channel :!relay_url bob https://www.youtube.com/watch?v=2XzmNpacpvk Interesting video"

	messageSlice := strings.Split(testMessage, " ")

	fromUser = strings.Split(messageSlice[0], ":")[1]
  fromUser = strings.Split(messageSlice[0], "!")[0]
  fromUser = strings.TrimLeft(fromUser, ":")
	toUser = messageSlice[4]
	channel = strings.Split(messageSlice[2], "!")[0]
	url = strings.Split(messageSlice[5], " ")[0]
	description = strings.Join(messageSlice[6:], " ")
	description = strings.TrimRight(description, "\r\n")

	help = "Usage: !relay_url <user> <url> <description>"

	fmt.Println("From User: ", fromUser)
	fmt.Println("To User: ", toUser)
  fmt.Println("Channel: ", channel)
	fmt.Println("Description: ", description)
	fmt.Println("URL: ", url)
	fmt.Println("Help: ", help)
  
  // TODO
  // - Map the variables to our struct
  // - Use the struct to save the message to the database
  // - Remap the database tables to just use the one relay_messages table
  // - Implement the help message
  // - Implement monitoring of the channel for the user and send the message to the user
  // - Mark the message as sent


}
