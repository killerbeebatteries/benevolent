package main

import (
	"fmt"
	"log"
	"math/rand"
)

type Greeting struct {
	ID        int
	FirstWord string
	Body      string
}

func getGreetings(name string) (string, error) {
	var greetings []Greeting
	var greetingCount int
	var randomGreeting Greeting

	defaultGreeting := fmt.Sprintf("Hello %s", name)

	rows, err := DB.Query("SELECT id, first_word, body FROM greetings")

	if err != nil {
		return defaultGreeting, err
	}

	defer rows.Close()

	for rows.Next() {
		var greeting Greeting
		err := rows.Scan(&greeting.ID, &greeting.FirstWord, &greeting.Body)

		if err != nil {
			return defaultGreeting, err
		}
		greetings = append(greetings, greeting)
	}

	if len(greetings) == 0 {
		return defaultGreeting, nil
	} else {
		greetingCount = len(greetings)
		randomGreeting = greetings[rand.Intn(greetingCount)]
	}

	return fmt.Sprintf("%s, %s, %s", randomGreeting.FirstWord, name, randomGreeting.Body), nil
}

func main() {

	err := OpenDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDatabase()

	fmt.Println(getGreetings("kronk"))
}
