package services

import (
	"ai-jira-triage/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var teamToAccount = map[string]string{
	"Frontend":      "712020:d44e700a-a4ae-497f-8d8b-7cb296946ec9",
	"Backend":       "712020:d44e700a-a4ae-497f-8d8b-7cb296946ec9",
	"Backoffice":    "712020:d44e700a-a4ae-497f-8d8b-7cb296946ec9",
	"User API":      "712020:d44e700a-a4ae-497f-8d8b-7cb296946ec9",
	"Notifications": "712020:d44e700a-a4ae-497f-8d8b-7cb296946ec9",
}

func UpdateIssue(issueKey string, aiResp *models.AIResponse) error {

	url := os.Getenv("JIRA_BASE_URL") + "/rest/api/3/issue/" + issueKey

	//  Get assignee from mapping
	assignee := teamToAccount[aiResp.Component]

	fmt.Println("Component:", aiResp.Component)
	fmt.Println("Assigning to accountId:", assignee)

	// Build fields dynamically
	fields := map[string]interface{}{
		"priority": map[string]string{
			"name": aiResp.Priority, 
		},
		"labels": []string{"ai-triaged"},
	}

	// Assign ONLY if High priority
	if aiResp.Priority == "High" && assignee != "" {
		fields["assignee"] = map[string]string{
			"accountId": assignee,
		}
	} else {
		fmt.Println("Skipping assignment (not High priority or no mapping)")
	}

	body := map[string]interface{}{
		"fields": fields,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	email := os.Getenv("JIRA_EMAIL")
	token := os.Getenv("JIRA_API_TOKEN")

	auth := base64.StdEncoding.EncodeToString([]byte(email + ":" + token))

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Jira Response:", string(bodyBytes))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("jira update failed: %d | %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Println("Jira updated successfully:", issueKey)

	return nil
}