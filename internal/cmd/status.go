package cmd

import (
	"github.com/spf13/cobra"
	"github.com/abdelhak/gh-issue-status/internal/pkg/application/status"
)

func newStatusCommand() *cobra.Command {
	var statusCommand = &cobra.Command{
		Use:   "change [issue-number]",
		Short: "Change issue status in project board",
		Long:  "Change the status of an issue in a GitHub project board (To Do, In Progress, Done).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			factory := status.NewFactory(command, args)
			return factory.Run()
		},
	}

	statusCommand.Flags().StringP("status", "s", "", "New status for the issue (todo, in-progress, done)")
	statusCommand.Flags().StringP("project", "p", "", "Project number or name")
	statusCommand.Flags().BoolP("json", "j", false, "Output in JSON format")

	return statusCommand
}