package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {

	root := cobra.Command{
		Use:   "ecs-deploy [flags]",
		Short: "ecs-deploy tool",
		Long:  "ecs-deploy tool",
	}

	taskDefCommand := &cobra.Command{
		Use:          "task-def",
		Short:        "Generate and optionally register Task Definitions",
		Long:         "Generate and optionally register Task Definitions",
		RunE:         taskDefFunction,
		SilenceUsage: true,
	}
	initTaskDefFlags(taskDefCommand.Flags())
	root.AddCommand(taskDefCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\necs-deploy completion > /usr/local/etc/bash_completion.d/ecs-deploy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
