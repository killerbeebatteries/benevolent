package main

import (
  "strings"
)

func getHelp(feature string) ([]string, error) {

	var help []string

  switch feature {

  case "relay_url":
 	  help = append(help, "Usage: !relay_url <user> <url> <description>")
	  help = append(help, "Description: Will post your message to the channel the next time the target user is active.")

  case "weather":
	  help = append(help, "Usage: !weather <location>")
  }

	return help, nil
}

func getUserFromMessage(message string) string {
	messageSlice := strings.Split(message, " ")

  user := strings.Split(messageSlice[0], "!")[0]
	user = strings.TrimLeft(user, ":")
  return user 
}
