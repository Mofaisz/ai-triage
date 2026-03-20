package handlers

import (
	"ai-jira-triage/models"
	"ai-jira-triage/services"
	"fmt"
	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {

	var payload models.WebhookPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		return
	}

	issue := payload.Issue

	fmt.Println("Webhook received:", issue.Key)

	aiResp, err := services.AnalyzeBug(issue.Fields.Summary, issue.Fields.Description)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("AI:", aiResp)

	// Confidence check
	if aiResp.Confidence < 0.8 {
		services.UpdateIssue(issue.Key, aiResp)
		c.JSON(200, gin.H{"status": "needs review"})
		return
	}

	// Update Jira
	err = services.UpdateIssue(issue.Key, aiResp)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Slack alert only for High
	if aiResp.Priority == "High" {
		services.SendSlackAlert(issue.Key, issue.Fields.Summary, aiResp.Priority)
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"priority": aiResp.Priority,
	})
}