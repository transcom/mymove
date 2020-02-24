package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {

	root := cobra.Command{
		Use:   "prime-api-client [flags]",
		Short: "Prime API client",
		Long:  "Prime API client",
	}

	fetchMTOsCommand := &cobra.Command{
		Use:          "fetch-mtos",
		Short:        "fetch mtos",
		Long:         "fetch move task orders",
		RunE:         fetchMTOs,
		SilenceUsage: true,
	}
	initFetchMTOsFlags(fetchMTOsCommand.Flags())
	root.AddCommand(fetchMTOsCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nprime-api-client completion > /usr/local/etc/bash_completion.d/prime-api-client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}