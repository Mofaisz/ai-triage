package services

import (
	"ai-jira-triage/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func AnalyzeBug(title string, description string) (*models.AIResponse, error) {

	apiKey := os.Getenv("OPENROUTER_API_KEY")

	url := "https://openrouter.ai/api/v1/chat/completions"

	if title == "" {
		title = "No title provided"
	}
	if description == "" {
		description = "No description provided"
	}
// prompt
	prompt := fmt.Sprintf(`
You are a Jira bug triage assistant.

STRICT RULES:
- Return ONLY valid JSON
- Do NOT include explanation
- Do NOT include markdown
- Always include all fields

{
  "priority": "High | Medium | Low",
  "component": "Frontend | Backend | Backoffice | User API | Notifications",
  "team": "Team Name",
  "confidence": number (0 to 1)
}

Title: %s
Description: %s
`, title, description)

	body := map[string]interface{}{
		"model": "openrouter/auto",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// ✅ Headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// ✅ Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// ✅ Parse response
	var resultMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&resultMap)

	// ✅ Extract AI text (this is tricky part)
	choices, ok := resultMap["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid AI response: %v", resultMap)
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format: %v", resultMap)
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format: %v", resultMap)
	}

	text := content
	fmt.Println("FULL AI RESPONSE:", resultMap)

	// ✅ Clean response
	text = strings.ReplaceAll(text, "```json", "")
	text = strings.ReplaceAll(text, "```", "")
	text = strings.TrimSpace(text)

	fmt.Println("AI RAW:", text)

	// ✅ Convert to struct
	var result models.AIResponse
	err = json.Unmarshal([]byte(text), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON parse error: %v | raw: %s", err, text)
	}

	return &result, nil
}
