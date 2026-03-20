package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	// "os"
)

func SendSlackAlert(issueKey, summary, priority string) {

	webhookURL := "***REMOVED***"
	message := fmt.Sprintf(
		"*Critical Bug Detected*\n\n*Issue:* %s\n*Summary:* %s\n*Priority:* %s\n*Link:* https://domindo.atlassian.net/browse/%s",
		issueKey, summary, priority, issueKey,
	)

	payload := map[string]string{
		"text": message,
	}

	jsonBody, _ := json.Marshal(payload)

	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
}
