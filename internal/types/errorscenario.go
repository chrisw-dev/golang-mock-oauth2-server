package types

// ErrorScenario defines an error scenario configuration
type ErrorScenario struct {
	Enabled     bool
	Endpoint    string
	StatusCode  int
	ErrorCode   string
	Description string
}
