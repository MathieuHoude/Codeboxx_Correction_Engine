package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GradingRequest contains the necessary elements to grade a project
type GradingRequest struct {
	GithubHandle string
}

//GradingResponse contains the informations to send back to the requester
type GradingResponse struct {
	testResults string
}

//RspecResults represents the data returned by Rspec
type RspecResults struct {
	Version  string `json:"version"`
	Examples []struct {
		ID              string      `json:"id"`
		Description     string      `json:"description"`
		FullDescription string      `json:"full_description"`
		Status          string      `json:"status"`
		FilePath        string      `json:"file_path"`
		LineNumber      int         `json:"line_number"`
		RunTime         float64     `json:"run_time"`
		PendingMessage  interface{} `json:"pending_message"`
	} `json:"examples"`
	Summary struct {
		Duration                     float64 `json:"duration"`
		ExampleCount                 int     `json:"example_count"`
		FailureCount                 int     `json:"failure_count"`
		PendingCount                 int     `json:"pending_count"`
		ErrorsOutsideOfExamplesCount int     `json:"errors_outside_of_examples_count"`
	} `json:"summary"`
	SummaryLine string `json:"summary_line"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Codeboxx Grading API")
	fmt.Println("Endpoint Hit: homePage")
}

func rubyResidentialControllerCorrection(w http.ResponseWriter, r *http.Request) {
	imageName := "mathieuhoude/ruby-residential-controller-grading"
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
	testResults := docker(imageName, request.GithubHandle)

	jsonData, err := json.Marshal(testResults)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/rubyresidentialcontroller", rubyResidentialControllerCorrection).Methods("POST")
	log.Println("Starting server on :10000...")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	loadEnv()
	handleRequests()
}
