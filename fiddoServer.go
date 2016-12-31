//fiddoServer.go
package main

import (
	"net/http"
	"log"
	"encoding/json"
)

type Commit struct {
	Id, Url string
}

type repo struct {
	Id, Name string
}

type body_struct struct {
	Ref string
	Commits []Commit
	Repository repo
}

func main() {
	http.HandleFunc("/webhook", webhook)
	http.ListenAndServe(":5000", nil)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body body_struct
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Printf("Ref be like: %v", body.Ref)
	log.Printf("Commits be like: %v", body.Commits)
	log.Printf("Repo be like: %v", body.Repository)
}