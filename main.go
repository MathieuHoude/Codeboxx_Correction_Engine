package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

// GradingRequest contains the necessary elements to grade a project
type GradingRequest struct {
	GithubHandle string
}

//GradingResponse contains the informations to send back to the requester
type GradingResponse struct {
	testResults map[string]interface{}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Codeboxx Grading API")
	fmt.Println("Endpoint Hit: homePage")
}

func dockerComposeBuild(githubHandle string) {

	cmd := exec.Command("docker-compose", "-f", "./rubyResidentialControllerGrading/docker-compose.yml", "--project-directory", "./rubyResidentialControllerGrading", "build", "--build-arg", "githubHandle="+githubHandle)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func dockerComposeUp() {
	cmd := exec.Command("docker-compose", "-f", "./rubyResidentialControllerGrading/docker-compose.yml", "--project-directory", "./rubyResidentialControllerGrading", "up")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func dockerComposeDown() {
	cmd := exec.Command("docker-compose", "-f", "./rubyResidentialControllerGrading/docker-compose.yml", "--project-directory", "./rubyResidentialControllerGrading", "down")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func dockerRun() []byte {
	cmd := exec.Command("docker", "run", "ruby-residential-controller-grading")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(string(stdout))

	return stdout
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
	docker(imageName, request.GithubHandle)

	// response := GradingResponse{
	// 	testResults: raw,
	// }

	// fmt.Printf("%+v\n", response)

	// jsonData, err := json.Marshal(response)

	// w.Header().Set("Content-Type", "application/json")
	// // w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(response)

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
