package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/abdelhak/gh-issue-status/internal/pkg/utils/cmdutil"
)

const version = "v1.0.0"

func Execute() {
	var rootCommand = &cobra.Command{
		Use:           "issue-status",
		Short:         "Manage GitHub issue status in project boards",
		Long:          "Change issue status in GitHub project boards (To Do, In Progress, Done).",
		SilenceErrors: true,
		RunE: func(command *cobra.Command, args []string) error {
			versionFlag, _ := command.Flags().GetBool("version")

			if versionFlag {
				fmt.Println(version)
				return nil
			}

			return command.Help()
		},
	}

	rootCommand.SetHelpFunc(cmdutil.HelpFunction)
	rootCommand.SetUsageFunc(cmdutil.UsageFunction)

	rootCommand.Flags().BoolP("version", "v", false, "Print the version of this extension")
	rootCommand.PersistentFlags().StringP("repo", "R", "", "Select another repository using the [HOST/]OWNER/REPO format")

	rootCommand.AddCommand(newStatusCommand())

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}