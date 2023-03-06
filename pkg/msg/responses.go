package msg

import "encoding/json"

// Result from calling the Lambda function.
type Result struct {
	// Response must be decoded into the expected response type
	// with json.Unmarshal().
	Response json.RawMessage `json:"response"`
}

type Resource struct {
	Type string         `json:"type"`
	ID   string         `json:"id"`
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

type LoadResponse struct {
	Resources []Resource    `json:"resources"`
	Tasks     []PendingTask `json:"tasks"`
}

type PendingTask struct {
	Task string         `json:"task"`
	Ctx  map[string]any `json:"ctx"`
}

type GrantResponse struct {
	AccessInstructions string         `json:"access_instructions"`
	State              map[string]any `json:"state"`
}
