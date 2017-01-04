//fiddoServer.go
package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"os/exec"
)

type HeadCommit struct {
	Id, Url string
}

type repo struct {
	Id int64
	Name string
}

type body_struct struct {
	Ref string				`json:"ref"`
	HeadCommit HeadCommit 	`json:"head_commit"`
	Repo repo 				`json:"repository"`
}

var scriptRoot = "/var/fiddoScripts/"

func main() {
	if (len(os.Args) > 1 && os.Args[1] != "") {
		scriptRoot = os.Args[1]
	}
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

	fullPathScript := scriptRoot + body.Repo.Name + ".sh"

	go executeFiddoScript(fullPathScript, body.HeadCommit.Url)

	log.Printf("Ref be like: %v", body.Ref)
	log.Printf("Commit be like: %v", body.HeadCommit)
	log.Printf("Repo be like: %v", body.Repo)
}

func executeFiddoScript(scriptName string, url string) {
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		log.Printf("Script not found: %v", scriptName)
		return
	}
	err := exec.Command(scriptName, url).Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		log.Printf("Script exited with error: %v", exitErr.Stderr)
		return
	} else if err != nil {
		panic(err)
	}
}