package release_branch

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	//"strconv"
	"strings"
)

func LatestTagFetch(typeOfTag string, header string) string {

	currentTagHeader := setupConfig.current_week_release
	// Define a dry-run flag
	//dryRun := flag.Bool("dryrun", false, "Run the program in dry-run mode without creating a release")
	//flag.Parse()

	// Fetch the latest tags from the remote repository
	err := fetchTags()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	tagPattern := ""

	// Define the tag pattern to match
	if typeOfTag == "RC" {
		tagPattern = currentTagHeader + header + "*-RC.*"
	} else if typeOfTag == "DEV" {
		tagPattern = currentTagHeader + header + "*-DEV.*"
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
		os.Exit(0)
	}

	// The first tag in the sorted output is the latest
	latestTag := tags[0]
	fmt.Println("Latest tag:", latestTag)

	return latestTag

	//// Increment the latest tag
	//newTag, err := incrementTag(latestTag)
	//if err != nil {
	//	fmt.Println("Error incrementing tag:", err)
	//	os.Exit(1)
	//}
	//fmt.Println("New tag:", newTag)
	//fmt.Printf("::set-output name=newTag::%typeOfTag\n", newTag)
	//
	//// Check if dry-run mode is enabled
	//if *dryRun {
	//	fmt.Println("Dry-run mode enabled. Skipping release creation.")
	//	fmt.Printf("Simulated release creation with tag: %typeOfTag\n", newTag)
	//	return ""
	//}
	//return newTag
}

// fetchTags fetches the latest tags from the remote repository
func fetchTags() error {
	cmd := exec.Command("git", "fetch", "--tags")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %s, output: %s", err, string(output))
	}
	return nil
}

// incrementTag increments the last numeric part of the tag by 1
func incrementTag(tag string) (string, error) {
	// Split the tag into parts by "."
	parts := strings.Split(tag, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid tag format: %s", tag)
	}

	// Parse the last part as an integer
	lastPart := parts[len(parts)-1]
	lastNum, err := strconv.Atoi(lastPart)
	if err != nil {
		return "", fmt.Errorf("failed to parse last part of tag as number: %s", lastPart)
	}

	// Increment the last number
	lastNum++

	// Reconstruct the tag with the incremented number
	parts[len(parts)-1] = strconv.Itoa(lastNum)
	return strings.Join(parts, "."), nil
}
