package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func versionFunction(cmd *cobra.Command, args []string) error {
	str, err := json.Marshal(map[string]interface{}{
		"gitBranch": gitBranch,
		"gitCommit": gitCommit,
	})
	if err != nil {
		return err
	}
	fmt.Println(string(str))
	return nil
}
