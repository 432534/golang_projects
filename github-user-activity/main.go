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
	fmt.Println("We are going to do a project where try to get details of commits , pushed done by the user in github")
	if len(os.Args) < 2 {
		fmt.Println("Provide the argument to help you")
		os.Exit(1)
	}
	username := os.Args[1]
	apiurl := "https://api.github.com/users/" + username + "/events"
	reponse, err := http.Get(apiurl)
	if err != nil {
		panic(err)
	}
	defer reponse.Body.Close()
	if reponse.StatusCode != http.StatusOK {
		fmt.Println("We are unable to fetch the details")
		os.Exit(1)
	}
	var events []Event
	decoder := json.NewDecoder(reponse.Body)
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
