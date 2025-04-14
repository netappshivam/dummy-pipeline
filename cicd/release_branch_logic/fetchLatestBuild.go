package release_branch_logic

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func LatestTagFetch(typeOfTag string) string {

	currentTagHeader := setupConfig.CurrentWeekRelease

	// Fetch the latest tags from the remote repository
	err := fetchTags()
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
	latestTag := tags[0]

	// Check if any tags were found
	if len(tags) == 0 || tags[0] == "" {
		fmt.Println("No tags found matching the pattern:", tagPattern)
		if typeOfTag == "DEV" {
			latestTag = currentTagHeader + ".0.0-DEV.1"
		}

	}

	// The first tag in the sorted output is the latest

	fmt.Println("Latest tag:", latestTag)

	return latestTag

}
