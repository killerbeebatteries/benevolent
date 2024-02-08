package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

type relayMessage struct {
	Id          int
	Timestamp   time.Time
	FromUser    string
	ToUser      string
	Description string
	URL         string
}

func saveMessage(message relayMessage) error {
	_, err := DB.Exec("INSERT INTO relay_messages (timestamp, from_user, to_user, description, suggested_url) VALUES ($1, $2, $3, $4, $5)", message.Timestamp, message.FromUser, message.ToUser, message.Description, message.URL)

	if err != nil {
		return err
	}

	return nil

}

func isValidURL(inputURL string) bool {
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return false
	}

	// Check if the URL scheme is present and is either "http" or "https"
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return true
}

func getHelp() []string {

	var help []string

	help = append(help, "Usage: !relay_url <user> <url> <description>")
	help = append(help, "Description: Will post your message to the channel the next time the target user is active.")

	return help
}

func relayUrlMessage(message string) ([]string, error) {
	err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDatabase()

	messageSlice := strings.Split(message, " ")

	help := getHelp()

	var fromUser, toUser, channel, description, url string
	var response []string

	fromUser = strings.Split(messageSlice[0], "!")[0]
	fromUser = strings.TrimLeft(fromUser, ":")
	toUser = messageSlice[4]
	channel = strings.Split(messageSlice[2], "!")[0]
	url = strings.Split(messageSlice[5], " ")[0]
	description = strings.Join(messageSlice[6:], " ")
	description = strings.TrimRight(description, "\r\n")

	if !isValidURL(url) {
		return help, fmt.Errorf("URL appears to be invalid: %s", url)
	}

	record := relayMessage{
		Timestamp:   time.Now(),
		FromUser:    fromUser,
		ToUser:      toUser,
		Description: description,
		URL:         url,
	}

	fmt.Println("Saving message to database")
	if err := saveMessage(record); err != nil {
		fmt.Println("Error saving message to database: ", err)
		return nil, err
	}

	fmt.Println("From User: ", fromUser)
	fmt.Println("To User: ", toUser)
	fmt.Println("Channel: ", channel)
	fmt.Println("Description: ", description)
	fmt.Println("URL: ", url)
	fmt.Println("Help: ", help)

	return response, nil

}

func main() {

	testMessage := ":alice!unknown@localhost PRIVMSG #channel :!relay_url bob https://www.youtube.com/watch?v=2XzmNpacpvk Interesting video"

	response, err := relayUrlMessage(testMessage)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(response)

	// TODO
	// - Implement the help message
	// - Implement monitoring of the channel for the user and send the message to the user
	// - Mark the message as sent

}
