package main

func startGrading(gradingRequest GradingRequest) GradingResponse {
	var gradingResponse GradingResponse
	switch gradingRequest.ProjectName {
	case "ruby-residential-controller":
		gradingResponse = rubyResidentialControllerCorrection(gradingRequest)

	}
	return gradingResponse

}

func rubyResidentialControllerCorrection(gradingRequest GradingRequest) GradingResponse {
	imageName := "mathieuhoude/ruby-residential-controller-grading"

	testResults := docker(imageName, gradingRequest.GithubHandle)
	issues := codeClimate(gradingRequest.GithubHandle, "Rocket_Elevators_Controllers")

	gradingResponse := GradingResponse{TestResults: testResults, Issues: issues}

	return gradingResponse

	// jsonData, err := json.Marshal(gradingResponse)

	// fmt.Println(string(jsonData))

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// w.Write(jsonData)

}
