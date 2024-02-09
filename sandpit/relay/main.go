package main

import (
	"errors"
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

func saveRelayMessage(message relayMessage) error {
	_, err := DB.Exec("INSERT INTO relay_messages (timestamp, from_user, to_user, description, suggested_url) VALUES ($1, $2, $3, $4, $5)", message.Timestamp, message.FromUser, message.ToUser, message.Description, message.URL)

	if err != nil {
		return err
	}

	return nil

}

func markRelayMessageAsSent(id int) error {
  _, err := DB.Exec("UPDATE relay_messages SET was_relayed = true WHERE id = $1", id)

  if err != nil {
    return err
  }

  return nil
}

func getRelayMessages() ([]relayMessage, error) {
  var messages []relayMessage

  rows, err := DB.Query("SELECT id, timestamp, from_user, to_user, description, suggested_url FROM relay_messages WHERE was_relayed = false)

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  for rows.Next() {
    var message relayMessage
    err := rows.Scan(&message.Id, &message.Timestamp, &message.FromUser, &message.ToUser, &message.Description, &message.URL)

    if err != nil {
      return nil, err
    }

    messages = append(messages, message)
  }

  return messages, nil
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

func isUserInChannel(user string, channel string, message string) bool {
  return true
}


func getHelp() []string {

	var help []string

	help = append(help, "Usage: !relay_url <user> <url> <description>")
	help = append(help, "Description: Will post your message to the channel the next time the target user is active.")

	return help
}

func addRelayMessage(message string) ([]string, error) {

	var fromUser, toUser, channel, description, url string
	var response []string

	messageSlice := strings.Split(message, " ")

	help := getHelp()

	fromUser = strings.Split(messageSlice[0], "!")[0]
	fromUser = strings.TrimLeft(fromUser, ":")
	toUser = messageSlice[4]
	channel = strings.Split(messageSlice[2], "!")[0]
	url = strings.Split(messageSlice[5], " ")[0]
	description = strings.Join(messageSlice[6:], " ")
	description = strings.TrimRight(description, "\r\n")

	if !isValidURL(url) {
    response = append(response, "URL appears to be invalid: " + url)
		return response, nil
	}

  if isUserInChannel(toUser, channel, message) {
    response = append(response, fmt.Sprintf("User %s is in the channel. Maybe they could just read this message? :D", toUser))
    return response, nil
  }

	record := relayMessage{
		Timestamp:   time.Now(),
		FromUser:    fromUser,
		ToUser:      toUser,
		Description: description,
		URL:         url,
	}

	err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDatabase()

	fmt.Println("Saving message to database")
	if err := saveRelayMessage(record); err != nil {
		fmt.Println("Error saving message to database: ", err)
    response = append(response, "Error saving message to database")
		return response, err
	}

	fmt.Println("From User: ", fromUser)
	fmt.Println("To User: ", toUser)
	fmt.Println("Channel: ", channel)
	fmt.Println("Description: ", description)
	fmt.Println("URL: ", url)
	fmt.Println("Help: ", help)

	return response, nil

}

func sendRelayMessage(toUser string) ([]string, error) {
  if messages, err := getRelayMessages(); err != nil {
    return "Error retrieving messages.", err
  }

  var response []string

  for _, message := range messages {
    if message.ToUser == toUser {
      response = append(response, fmt.Sprintf("%s: %s %s", message.FromUser, message.Description, message.URL))
      if err := markRelayMessageAsSent(message.Id); err != nil {
        return "Error marking message as sent.", err
      }
    }  
  } 

  if len(response) == 0 {
    response = append(response, fmt.Sprint("Hello %s. I have no pending messages for you.", toUser))
  }

  return response, nil

}

func main() {

	testMessage := ":alice!unknown@localhost PRIVMSG #channel :!relay_url bob https://www.youtube.com/watch?v=2XzmNpacpvk Interesting video"

	response, err := addRelayMessage(testMessage)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(response)

	// TODO
	// - Implement monitoring of the channel for the user and send the message to the user
	// - Mark the message as sent

}
