package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "Health Checker",
	Short: "api health checker application based on clean architecture",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
