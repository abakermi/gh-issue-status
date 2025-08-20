package status

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/abdelhak/gh-issue-status/internal/pkg/infrastructure/gh"
)

type Factory struct {
	command *cobra.Command
	args    []string
}

func NewFactory(command *cobra.Command, args []string) *Factory {
	return &Factory{
		command: command,
		args:    args,
	}
}

func (f *Factory) Run() error {
	repo, _ := f.command.Flags().GetString("repo")
	status, _ := f.command.Flags().GetString("status")
	project, _ := f.command.Flags().GetString("project")
	jsonOutput, _ := f.command.Flags().GetBool("json")

	var issueNumber int
	var err error

	if len(f.args) > 0 {
		issueNumber, err = strconv.Atoi(f.args[0])
		if err != nil {
			return fmt.Errorf("invalid issue number: %s", f.args[0])
		}
	} else {
		// List issues using gh CLI
		args := []string{"issue", "list", "--json", "number,title"}
		if repo != "" {
			args = append(args, "--repo", repo)
		}
		
		result, err := gh.Execute(args)
		if err != nil {
			return fmt.Errorf("failed to list issues: %v", err)
		}

		var issues []struct {
			Number int    `json:"number"`
			Title  string `json:"title"`
		}
		
		if err := json.Unmarshal([]byte(result), &issues); err != nil {
			return fmt.Errorf("failed to parse issues: %v", err)
		}

		var issueOptions []string
		for _, issue := range issues {
			issueOptions = append(issueOptions, fmt.Sprintf("#%d - %s", issue.Number, issue.Title))
		}

		var selected string
		prompt := &survey.Select{
			Message: "Select an issue:",
			Options: issueOptions,
		}
		survey.AskOne(prompt, &selected)
		fmt.Sscanf(selected, "#%d", &issueNumber)
	}

	if status == "" {
		statusOptions := []string{"todo", "in-progress", "done"}
		prompt := &survey.Select{
			Message: "Select new status:",
			Options: statusOptions,
		}
		survey.AskOne(prompt, &status)
	}

	if project == "" {
		repoInfo, err := gh.RetrieveRepoInformation(repo)
		if err != nil {
			return err
		}
		
		args := []string{"project", "list", "--owner", repoInfo.Owner, "--format", "json"}
		result, err := gh.Execute(args)
		if err != nil {
			return fmt.Errorf("failed to list projects: %v", err)
		}

		var response struct {
			Projects []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			} `json:"projects"`
		}
		
		if err := json.Unmarshal([]byte(result), &response); err != nil {
			return fmt.Errorf("failed to parse projects: %v", err)
		}

		if len(response.Projects) == 0 {
			return fmt.Errorf("no projects found for %s", repoInfo.Owner)
		}

		var projectOptions []string
		for _, proj := range response.Projects {
			projectOptions = append(projectOptions, proj.Title)
		}

		prompt := &survey.Select{
			Message: "Select project:",
			Options: projectOptions,
		}
		survey.AskOne(prompt, &project)
	}

	// First, add issue to project if not already added
	args := []string{"issue", "edit", fmt.Sprintf("%d", issueNumber), "--add-project", project}
	if repo != "" {
		args = append(args, "--repo", repo)
	}
	
	_, err = gh.Execute(args)
	if err != nil {
		// Issue might already be in project, continue
		fmt.Printf("Note: Issue may already be in project\n")
	}

	// Update the issue status in the project with spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Updating issue status..."
	s.Start()
	
	err = gh.UpdateIssueStatus(repo, issueNumber, project, status)
	s.Stop()
	
	if err != nil {
		return fmt.Errorf("failed to update issue status: %v", err)
	}

	if jsonOutput {
		result := map[string]interface{}{
			"issue":   issueNumber,
			"status":  status,
			"project": project,
		}
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	} else {
		fmt.Printf("✓ Issue #%d status changed to '%s' in project '%s'\n", issueNumber, status, project)
	}

	return nil
}