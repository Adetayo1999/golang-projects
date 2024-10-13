package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type GithubUserEvent struct {
	ID string `json:"id"`
    Type string `json:"type"`
    Actor struct {
      ID int `json:"id"`
      Login string `json:"login"`
      DisplayLogin string `json:"display_login"`
    } `json:"actor"`
    Repo  struct{
      ID int `json:"id"`
      Name string `json:"name"`
      URL string `json:"url"`
    } `json:"repo"`
	Payload struct {
	  RepositoryID *int `json:"repository_id"`
      PushID int `json:"push_id"`
	  Commits *[]interface{} `json:"commits"`
	  Ref *string `json:"ref"`
      RefType *string `json:"ref_type"`
      MasterBranch *string `json:"master_branch"`
	} `json:"payload"`
	Public bool `json:"public"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {

	var events []GithubUserEvent

	username := flag.String("username", "Adetayo1999", "i.e username=Adetayo1999")

	flag.Parse()

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", *username))

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.New("something went wrong"))
	}

	err = json.NewDecoder(resp.Body).Decode(&events);

	if err != nil {
		panic(err);
	}

	for _, event := range events {

		switch{
		case event.Type == "PushEvent":
			fmt.Printf("- Pushed %d commits to %s on %s \n", len(*event.Payload.Commits), event.Repo.Name, event.CreatedAt.Format(time.RFC850))
		
		case event.Type == "CreateEvent":
			fmt.Printf("- Created %s on repository %s on %s \n", *event.Payload.RefType, event.Repo.Name, event.CreatedAt.Format(time.RFC850))	
			
		case event.Type == "WatchEvent":
			fmt.Printf("- %v starred %s(%s) on %s \n", event.Actor.DisplayLogin, event.Repo.Name, event.Repo.URL, event.CreatedAt.Format(time.RFC850))	
		}

	}

}
