package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload map[string]any `json:"payload"`
}

func main() {
	fmt.Println("We are building CLI for github-user-activity")

	if len(os.Args) < 2 {
		fmt.Println("Provide a valid GitHub username as an argument")
		os.Exit(1)
	}

	username := os.Args[1]
	apiurl := "https://api.github.com/users/" + username + "/events"

	response, err := http.Get(apiurl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch data from API: %s (status %d)\n", response.Status, response.StatusCode)
		os.Exit(1)
	}

	var events []Event
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&events)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No recent public activity found")
		return
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commits, ok := event.Payload["commits"].([]interface{})
			if ok {
				fmt.Printf("Pushed %d commits to %s\n", len(commits), event.Repo.Name)
			} else {
				fmt.Printf("Pushed to %s\n", event.Repo.Name)
			}

		case "IssuesEvent":
			action, ok := event.Payload["action"].(string)
			if ok {
				fmt.Printf("%s an issue in %s\n", action, event.Repo.Name)
			}

		case "IssueCommentEvent":
			fmt.Printf("Commented on an issue in %s\n", event.Repo.Name)

		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)

		case "ForkEvent":
			fmt.Printf("Forked %s\n", event.Repo.Name)

		default:
			fmt.Printf("%s in %s\n", event.Type, event.Repo.Name)
		}
	}
}
