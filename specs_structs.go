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

//JestResults contains the data returned by Jest
type JestResults struct {
	NumFailedTestSuites       int           `json:"numFailedTestSuites"`
	NumFailedTests            int           `json:"numFailedTests"`
	NumPassedTestSuites       int           `json:"numPassedTestSuites"`
	NumPassedTests            int           `json:"numPassedTests"`
	NumPendingTestSuites      int           `json:"numPendingTestSuites"`
	NumPendingTests           int           `json:"numPendingTests"`
	NumRuntimeErrorTestSuites int           `json:"numRuntimeErrorTestSuites"`
	NumTodoTests              int           `json:"numTodoTests"`
	NumTotalTestSuites        int           `json:"numTotalTestSuites"`
	NumTotalTests             int           `json:"numTotalTests"`
	OpenHandles               []interface{} `json:"openHandles"`
	Snapshot                  struct {
		Added               int           `json:"added"`
		DidUpdate           bool          `json:"didUpdate"`
		Failure             bool          `json:"failure"`
		FilesAdded          int           `json:"filesAdded"`
		FilesRemoved        int           `json:"filesRemoved"`
		FilesRemovedList    []interface{} `json:"filesRemovedList"`
		FilesUnmatched      int           `json:"filesUnmatched"`
		FilesUpdated        int           `json:"filesUpdated"`
		Matched             int           `json:"matched"`
		Total               int           `json:"total"`
		Unchecked           int           `json:"unchecked"`
		UncheckedKeysByFile []interface{} `json:"uncheckedKeysByFile"`
		Unmatched           int           `json:"unmatched"`
		Updated             int           `json:"updated"`
	} `json:"snapshot"`
	StartTime   int64 `json:"startTime"`
	Success     bool  `json:"success"`
	TestResults []struct {
		AssertionResults []struct {
			AncestorTitles  []string      `json:"ancestorTitles"`
			FailureMessages []interface{} `json:"failureMessages"`
			FullName        string        `json:"fullName"`
			Location        interface{}   `json:"location"`
			Status          string        `json:"status"`
			Title           string        `json:"title"`
		} `json:"assertionResults"`
		EndTime   int64  `json:"endTime"`
		Message   string `json:"message"`
		Name      string `json:"name"`
		StartTime int64  `json:"startTime"`
		Status    string `json:"status"`
		Summary   string `json:"summary"`
	} `json:"testResults"`
	WasInterrupted bool `json:"wasInterrupted"`
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