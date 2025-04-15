/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dummy-pipeline/cmd/jira"
	"dummy-pipeline/cmd/tag"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vsactl",
	Short: "A cli used to control vsa cicd controller",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	_ = godotenv.Load()
	rootCmd.AddCommand(jira.JiraCmd)
	rootCmd.AddCommand(tag.TagCmd)
}
