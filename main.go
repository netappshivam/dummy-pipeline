package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Fetch all tags from the remote
	cmd := exec.Command("git", "fetch", "--tags")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	// Get all tag names
	cmd = exec.Command("git", "tag")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error listing tags:", err)
		os.Exit(1)
	}
	rawTags := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter and sort the tags
	tags := make([]string, 0, len(rawTags))
	for _, tag := range rawTags {
		if strings.HasPrefix(tag, "v") {
			tags = append(tags, tag)
		}
	}
	sort.Strings(tags)

	// Find the latest tag
	if len(tags) == 0 {
		fmt.Println("No tags found, creating the first tag v0.0.1")
		fmt.Println("::set-output name=new_tag::v0.0.1")
		os.Exit(0)
	}
	latestTag := tags[len(tags)-1]

	// Increment the tag
	parts := strings.Split(strings.TrimPrefix(latestTag, "v"), ".")
	if len(parts) != 3 {
		fmt.Println("Latest tag does not follow semantic versioning:", latestTag)
		os.Exit(1)
	}
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	patch++
	newTag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)

	// Output the new tag to be used by other steps
	fmt.Printf("::set-output name=new_tag::%s\n", newTag)
}
