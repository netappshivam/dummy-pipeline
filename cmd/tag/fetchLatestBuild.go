package tag

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var FetchLatestBuildCmd = &cobra.Command{
	Use:   "fetch_latest_build",
	Short: "Fetch the latest build tag based on the type of tag (RC or DEV)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Println("Error: Please provide the type of tag (RC or DEV).")
			os.Exit(1)
		}
		typeOfTag := args[0]
		latestTag := LatestTagFetch(typeOfTag)
		fmt.Printf("Latest tag for %s: %s\n", typeOfTag, latestTag)
		return nil
	},
}

func LatestTagFetch(typeOfTag string) string {

	currentTagHeader := SetupConfigobject.CurrentWeekRelease

	// Fetch the latest tags from the remote repository
	err := FetchTags()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	tagPattern := ""

	// Define the tag pattern to match
	if typeOfTag == "RC" {
		tagPattern = currentTagHeader + "*-RC.*"
	} else if typeOfTag == "DEV" {
		tagPattern = currentTagHeader + "*-DEV.*"
	}

	// Execute the git command to list tags sorted by version (descending)
	cmd := exec.Command("git", "tag", "-l", "--sort=-v:refname", tagPattern)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing git command:", err)
		os.Exit(1)
	}

	// Split the output into individual tags
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Check if any tags were found
	if len(tags) == 0 || tags[0] == "" {
		fmt.Println("No tags found matching the pattern:", tagPattern)
	}

	// The first tag in the sorted output is the latest
	latestTag := tags[0]
	fmt.Println("Latest tag:", latestTag)

	return latestTag

}

func init() {
	// Register the FetchLatestBuildCmd with the parent command
	TagCmd.AddCommand(FetchLatestBuildCmd)
}
