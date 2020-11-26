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
	TestResults RspecResults
	Issues      CodeClimateIssues
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
	issues := codeClimate(request.GithubHandle, "Rocket_Elevators_Controllers")

	gradingResponse := GradingResponse{TestResults: testResults, Issues: issues}

	jsonData, err := json.Marshal(gradingResponse)

	fmt.Println(string(jsonData))

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
