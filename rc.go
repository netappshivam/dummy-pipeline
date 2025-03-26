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
	cmd = exec.Command("git", "tag", "--list", "v*-RC.*")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error listing tags:", err)
		os.Exit(1)
	}
	rawTags := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter tags that follow semantic versioning
	semverTags := make([]string, 0)
	for _, tag := range rawTags {
		if strings.HasPrefix(tag, "v") {
			parts := strings.Split(strings.TrimPrefix(tag, "v"), "-")
			if len(parts) == 2 {
				semverTags = append(semverTags, tag)
			}
		}
	}

	// Sort tags using semantic versioning rules
	sort.Slice(semverTags, func(i, j int) bool {
		return versionCompare_RC(semverTags[i], semverTags[j])
	})

	// No tags found, start with initial version
	if len(semverTags) == 0 {
		fmt.Println("No tags found, creating the first tag v0.0.0-RC.1")
		fmt.Println("::set-output name=new_tag::v0.0.0-RC.1")
		os.Exit(0)
	}

	// Get the latest tag
	latestTag := semverTags[len(semverTags)-1]

	// Increment the tag
	parts := strings.Split(strings.TrimPrefix(latestTag, "v"), "-")
	if len(parts) != 2 {
		fmt.Println("Latest tag does not follow semantic versioning:", latestTag)
		os.Exit(1)
	}
	rcParts := strings.Split(parts[1], ".")
	if len(rcParts) != 2 || rcParts[0] != "RC" {
		fmt.Println("Latest tag does not follow semantic versioning:", latestTag)
		os.Exit(1)
	}
	rcNumber, _ := strconv.Atoi(rcParts[1])
	rcNumber++
	newTag := fmt.Sprintf("v%s-RC.%d", parts[0], rcNumber)

	// Output the new tag to be used by other steps
	fmt.Printf("::set-output name=new_tag::%s\n", newTag)
}

// versionCompare_RC compares two semantic versioning tags.
func versionCompare_RC(tag1, tag2 string) bool {
	version1 := strings.TrimPrefix(tag1, "v")
	version2 := strings.TrimPrefix(tag2, "v")

	v1Parts := strings.Split(version1, "-")
	v2Parts := strings.Split(version2, "-")

	if v1Parts[0] != v2Parts[0] {
		return v1Parts[0] < v2Parts[0]
	}

	rc1, _ := strconv.Atoi(strings.TrimPrefix(v1Parts[1], "RC."))
	rc2, _ := strconv.Atoi(strings.TrimPrefix(v2Parts[1], "RC."))
	return rc1 < rc2
}
