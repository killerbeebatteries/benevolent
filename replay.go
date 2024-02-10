package main

import (
	"fmt"
	"log"
	"strings"
	"time"
  "errors"
  "net/url"
)

type relayMessage struct {
	Id          int
	Timestamp   time.Time
	FromUser    string
  FromChannel string
	ToUser      string
	Description string
	URL         string
}

func saveRelayMessage(message relayMessage) error {
	_, err := DB.Exec("INSERT INTO relay_messages (timestamp, from_user, from_channel, to_user, description, suggested_url) VALUES ($1, $2, $3, $4, $5, $6)", message.Timestamp, message.FromUser, message.FromChannel, message.ToUser, message.Description, message.URL)

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

func getRelayMessageFromCommand(message string) (relayMessage, error) {

	var fromUser, toUser, channel, description, url string

	messageSlice := strings.Split(message, " ")

	fromUser = strings.Split(messageSlice[0], "!")[0]
	fromUser = strings.TrimLeft(fromUser, ":")
	fromUser = strings.ToLower(fromUser)
	toUser = strings.ToLower(messageSlice[4])
	channel = strings.Split(messageSlice[2], "!")[0]
	url = strings.Split(messageSlice[5], " ")[0]
	description = strings.Join(messageSlice[6:], " ")
	description = strings.TrimRight(description, "\r\n")

  // TODO: I am not happy with the lack of error checking when assigning the values to the struct. I should probably add some checks here.
	record := relayMessage{
		Timestamp:   time.Now(),
		FromUser:    fromUser,
    FromChannel: channel,
		ToUser:      toUser,
		Description: description,
		URL:         url,
	}
  
  return record, nil
}


func getRelayMessages() ([]relayMessage, error) {
	var messages []relayMessage

	rows, err := DB.Query("SELECT id, timestamp, from_user, to_user, description, suggested_url FROM relay_messages WHERE was_relayed = false")

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
	// TODO
	return false
}

func addRelayMessage(message string) ([]string, error) {

	var response []string

  record, err := getRelayMessageFromCommand(message)

  if err != nil {
    response = append(response, "Error parsing message")
    return response, err
  }

	if !isValidURL(record.URL) {
    response = append(response, "Invalid URL: " + record.URL)
    return response, errors.New("Invalid URL: " + record.URL)
  }
	
	if isUserInChannel(record.ToUser, record.FromChannel, message) {
		response = append(response, fmt.Sprintf("User %s is in the channel. Maybe they could just read this message? :D", record.ToUser))
		return response, nil
	}

	err = OpenDatabase()
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

  if len(response) == 0 {
    response = append(response, fmt.Sprintf("Message saved for %s. I will relay it the next time they are kicking around here.", record.ToUser))
  }
	return response, nil

}

func sendRelayMessage(toUser string) ([]string, error) {
	var response []string

	err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDatabase()

	messages, err := getRelayMessages()

	if err != nil {
		response = []string{"Error retrieving messages."}
		return response, err
	}

	for _, message := range messages {
		if message.ToUser == toUser {
			response = append(response, fmt.Sprintf("%s: %s %s", message.FromUser, message.Description, message.URL))
			if err := markRelayMessageAsSent(message.Id); err != nil {
				response = []string{"Error marking message as sent."}
				return response, err
			}
		}
	}

	if len(response) == 0 {
		response = append(response, fmt.Sprintf("Hello %s. I have no pending messages for you.", toUser))
	}

	return response, nil

}
