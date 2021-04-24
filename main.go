package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Codeboxx Grading API")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/correctionrequest", newCorrectionRequest).Methods("POST")
	myRouter.HandleFunc("/gradingrequest", newGradingRequest).Methods("POST")
	myRouter.HandleFunc("/checkGithubAccess", checkGithubAccessRequest).Methods("POST")

	log.Println("Starting server on :10000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:10000", myRouter))
}

func checkGithubAccessRequest(w http.ResponseWriter, r *http.Request) {

}

func newCorrectionRequest(w http.ResponseWriter, r *http.Request) {
	var request CorrectionRequest
	// request.DeliverableDeadline.Format("2020-12-05 13:31:50")
	err := decodeJSONBody(w, r, &request)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	updateJobStatus(request.JobID, "Received for correction")
	body, err := json.Marshal(request)
	if err == nil {
		newTask(request.JobID, body, "correction")
	}
	fmt.Fprintf(w, "The request has been received")
	fmt.Println("Endpoint Hit: correctionRequest")
}

func newGradingRequest(w http.ResponseWriter, r *http.Request) {
	var request GradingRequest
	// request.DeliverableDeadline.Format("2020-12-05 13:31:50")
	err := decodeJSONBody(w, r, &request)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	updateJobStatus(request.JobID, "Received for grading")
	body, err := json.Marshal(request)
	if err == nil {
		newTask(request.JobID, body, "grading")
	}
	fmt.Fprintf(w, "The request has been received")
	fmt.Println("Endpoint Hit: gradingrequest")
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://657e3fa075324ae2b5e3c4a81621bef9@o481104.ingest.sentry.io/5602523",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	loadEnv()
	startWorkers(5, 2) //Starts the workers that will receive tasks from the task_queue. Specify the number of workers needed.
	handleRequests()   //Start the API to accept and dispatch new grading requests
}
