package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GradingRequest contains the necessary elements to grade a project
type GradingRequest struct {
	JobID               int                `json:"JobID"`
	DeliverableID       int                `json:"DeliverableID"`
	DeliverableScores   []DeliverableScore `json:"DeliverableScores"`
	DeliverableDeadline time.Time          `json:"DeliverableDeadline"`
	GithubHandle        string             `json:"GithubHandle"`
	RepositoryURL       string             `json:"RepositoryURL"`
	DockerImageName     string             `json:"DockerImageName"`
	TestingTool         string             `json:"TestingTool"`
}

//GradingResponse contains the informations to send back to the requester
type GradingResponse struct {
	JobID             int
	DeliverableID     int
	DeliverableScores []DeliverableScore
	Issues            CodeClimateIssues
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
	fmt.Println(request.DeliverableDeadline)
	updateJobStatus(request.JobID, "Received")
	newTask(request)
	fmt.Fprintf(w, "The request has been received")
	fmt.Println("Endpoint Hit: gradingrequest")
}

func main() {
	loadEnv()
	// var x []DeliverableScore
	// getLastCommitDate(x, "MathieuHoude/Rocket_Elevators_Controllers")
	// codeClimate("MathieuHoude/Rocket_Elevators_Controllers")

	startWorkers(5)  //Starts the workers that will receive tasks from the task_queue. Specify the number of workers needed.
	handleRequests() //Start the API to accept and dispatch new grading requests
}
