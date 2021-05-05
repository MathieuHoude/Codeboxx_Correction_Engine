package main

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

//TestResults contains the data returned by Jest
type TestResult struct {
	ScoreCardItemName string `json:"ScoreCardItemName"`
	Status            string `json:"Status"`
}

//PytestResults contains the data returned by Pytest
type PytestResults struct {
	Created     float64 `json:"created"`
	Duration    float64 `json:"duration"`
	Exitcode    int     `json:"exitcode"`
	Root        string  `json:"root"`
	Environment struct {
		Python   string `json:"Python"`
		Platform string `json:"Platform"`
		Packages struct {
			Pytest string `json:"pytest"`
			Py     string `json:"py"`
			Pluggy string `json:"pluggy"`
		} `json:"Packages"`
		Plugins struct {
			Metadata   string `json:"metadata"`
			JSONReport string `json:"json-report"`
		} `json:"Plugins"`
	} `json:"environment"`
	Summary struct {
		Passed    int `json:"passed"`
		Total     int `json:"total"`
		Collected int `json:"collected"`
	} `json:"summary"`
	Tests []struct {
		Nodeid   string   `json:"nodeid"`
		Lineno   int      `json:"lineno"`
		Outcome  string   `json:"outcome"`
		Keywords []string `json:"keywords"`
		Setup    struct {
			Duration float64 `json:"duration"`
			Outcome  string  `json:"outcome"`
		} `json:"setup"`
		Call struct {
			Duration float64 `json:"duration"`
			Outcome  string  `json:"outcome"`
		} `json:"call"`
		Teardown struct {
			Duration float64 `json:"duration"`
			Outcome  string  `json:"outcome"`
		} `json:"teardown"`
	} `json:"tests"`
}
