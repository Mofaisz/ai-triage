package models

type WebhookPayload struct {
	Issue Issue `json:"issue"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Labels      []string `json:"labels"`
}

// AI response struct
type AIResponse struct {
	Priority   string  `json:"priority"`
	Component  string  `json:"component"`
	Team       string  `json:"team"`
	Confidence float64 `json:"confidence"`
}