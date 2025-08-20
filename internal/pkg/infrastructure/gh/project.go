package gh

import (
	"encoding/json"
	"fmt"
	"strings"
)

func UpdateIssueStatus(repo string, issueNumber int, projectTitle, status string) error {
	repoInfo, err := RetrieveRepoInformation(repo)
	if err != nil {
		return err
	}

	// Get project ID from project list
	args := []string{"project", "list", "--owner", repoInfo.Owner, "--format", "json"}
	result, err := Execute(args)
	if err != nil {
		return err
	}

	var response struct {
		Projects []struct {
			ID     string `json:"id"`
			Number int    `json:"number"`
			Title  string `json:"title"`
		} `json:"projects"`
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return err
	}

	var projectID string
	var projectNumber int
	for _, proj := range response.Projects {
		if proj.Title == projectTitle {
			projectID = proj.ID
			projectNumber = proj.Number
			break
		}
	}

	if projectID == "" {
		return fmt.Errorf("project '%s' not found", projectTitle)
	}

	// Get project items using project number
	args = []string{"project", "item-list", fmt.Sprintf("%d", projectNumber), "--owner", repoInfo.Owner, "--format", "json"}
	result, err = Execute(args)
	if err != nil {
		return fmt.Errorf("failed to list project items: %v", err)
	}

	var itemsResponse struct {
		Items []struct {
			ID      string `json:"id"`
			Content struct {
				Number int `json:"number"`
			} `json:"content"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(result), &itemsResponse); err != nil {
		return err
	}

	var itemID string
	for _, item := range itemsResponse.Items {
		if item.Content.Number == issueNumber {
			itemID = item.ID
			break
		}
	}

	if itemID == "" {
		// Issue not in project, add it first
		args = []string{"project", "item-add", fmt.Sprintf("%d", projectNumber), "--owner", repoInfo.Owner, "--url", fmt.Sprintf("https://github.com/%s/%s/issues/%d", repoInfo.Owner, repoInfo.Name, issueNumber)}
		_, err = Execute(args)
		if err != nil {
			return fmt.Errorf("failed to add issue to project: %v", err)
		}
		
		// Re-fetch project items to get the new item ID
		args = []string{"project", "item-list", fmt.Sprintf("%d", projectNumber), "--owner", repoInfo.Owner, "--format", "json"}
		result, err = Execute(args)
		if err != nil {
			return fmt.Errorf("failed to list project items after adding: %v", err)
		}
		
		if err := json.Unmarshal([]byte(result), &itemsResponse); err != nil {
			return err
		}
		
		for _, item := range itemsResponse.Items {
			if item.Content.Number == issueNumber {
				itemID = item.ID
				break
			}
		}
		
		if itemID == "" {
			return fmt.Errorf("failed to find issue after adding to project")
		}
	}

	// Get project fields to find Status field
	args = []string{"project", "field-list", fmt.Sprintf("%d", projectNumber), "--owner", repoInfo.Owner, "--format", "json"}
	result, err = Execute(args)
	if err != nil {
		return err
	}

	var fieldsResponse struct {
		Fields []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Options []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"options"`
		} `json:"fields"`
	}

	if err := json.Unmarshal([]byte(result), &fieldsResponse); err != nil {
		return err
	}

	var fieldID, optionID string
	statusMap := map[string]string{
		"todo":        "Todo",
		"in-progress": "In Progress", 
		"done":        "Done",
	}

	targetStatus := statusMap[status]
	if targetStatus == "" {
		targetStatus = status
	}

	for _, field := range fieldsResponse.Fields {
		if strings.ToLower(field.Name) == "status" {
			fieldID = field.ID
			for _, option := range field.Options {
				if strings.EqualFold(option.Name, targetStatus) {
					optionID = option.ID
					break
				}
			}
			break
		}
	}

	if fieldID == "" {
		return fmt.Errorf("status field not found in project")
	}
	if optionID == "" {
		return fmt.Errorf("status option '%s' not found", targetStatus)
	}

	// Update the item field value
	args = []string{"project", "item-edit", "--id", itemID, "--project-id", projectID, "--field-id", fieldID, "--single-select-option-id", optionID}
	_, err = Execute(args)
	return err
}