package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

func HelpFunction(command *cobra.Command, args []string) {
	fmt.Print(command.Long)

	if command.Runnable() {
		fmt.Printf("\n\nUSAGE\n  %s\n", command.UseLine())
	}

	if len(command.Commands()) > 0 {
		fmt.Printf("\n\nCOMMANDS\n")
		for _, c := range command.Commands() {
			fmt.Printf("  %s\n", rpad(c.Name(), c.NamePadding())+c.Short)
		}
	}

	if command.HasAvailableLocalFlags() {
		fmt.Printf("\n\nFLAGS\n")
		fmt.Print(command.LocalFlags().FlagUsages())
	}

	if command.HasAvailableInheritedFlags() {
		fmt.Printf("\n\nINHERITED FLAGS\n")
		fmt.Print(command.InheritedFlags().FlagUsages())
	}

	fmt.Printf("\n\nLEARN MORE\n")
	fmt.Printf("  Use 'gh issue-status <command> --help' for more information about a command.\n")
}

func UsageFunction(command *cobra.Command) error {
	command.Printf("Usage: %s\n", command.UseLine())
	return nil
}

func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds ", padding)
	return fmt.Sprintf(template, s)
}
