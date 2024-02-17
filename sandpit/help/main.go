package main

import (
  "strings"
  "fmt"
)

func getHelp(feature string) ([]string, error) {

	var help []string

  switch feature {

  case "time":
    help = append(help, "Usage: !time")
    help = append(help, "Description: Will return the current time.")

  case "ping":
    help = append(help, "Usage: !ping")
    help = append(help, "Description: Will return 'pong'.")

  case "hello":
    help = append(help, "Usage: !hello")
    help = append(help, "Description: Say hello.")

  case "relay_url":
 	  help = append(help, "Usage: !relay_url <user> <url> <description>")
	  help = append(help, "Description: Will post your message to the channel the next time the target user is active.")

  case "weather":
	  help = append(help, "Usage: !weather <location>")
    help = append(help, "Description: Will return the current weather for the specified location.")
  default:
    help = append(help, fmt.Sprintf("Sorry, feature %s was not found.", feature))
    help = append(help, "Available features are: time, ping, hello, relay_url, weather")
    help = append(help, "Usage: !help <feature>")
  }

	return help, nil
}

func getUserFromMessage(message string) string {
	messageSlice := strings.Split(message, " ")

  user := strings.Split(messageSlice[0], "!")[0]
	user = strings.TrimLeft(user, ":")
  return user 
}

func main() {
  getHelp("all")
}
