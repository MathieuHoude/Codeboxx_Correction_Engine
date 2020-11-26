package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GradingRequest contains the necessary elements to grade a project
type GradingRequest struct {
	GithubHandle  string
	RepositoryURL string
	ProjectName   string
}

//GradingResponse contains the informations to send back to the requester
type GradingResponse struct {
	TestResults RspecResults
	Issues      CodeClimateIssues
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Codeboxx Grading API")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// myRouter.HandleFunc("/ruby-residential-controller", rubyResidentialControllerCorrection).Methods("POST")
	myRouter.HandleFunc("/gradingrequest", newGradingRequest).Methods("POST")
	log.Println("Starting server on :10000...")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func newGradingRequest(w http.ResponseWriter, r *http.Request) {
	var request GradingRequest
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
	newTask(request)
	fmt.Fprintf(w, "The request has been received")
	fmt.Println("Endpoint Hit: gradingrequest")
}

func main() {
	loadEnv()
	startWorkers(3)
	handleRequests()
}
