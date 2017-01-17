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
	HeadCommit HeadCommit 	`json:"head_commit"`
	Repo repo 				`json:"repository"`
}

var scriptRoot = "/var/fiddoScripts/"

func main() {
	log.Printf("Fiddo listening on port: %d", 5000)
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

	go executeFiddoScript(fullPathScript, body.HeadCommit.Id)
}

func executeFiddoScript(scriptName string, id string) {
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		log.Printf("Script not found: %v", scriptName)
		return
	}
	log.Printf("Running this: %s %s", scriptName, id)
	err := exec.Command(scriptName, id).Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		log.Printf("Script exited with error: %v", exitErr)
		return
	} else if err != nil {
		panic(err)
	}
}
