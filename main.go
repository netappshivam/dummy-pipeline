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
	//
	//fmt.Printf("Found %d tags\n", len(rawTags))
	//
	//for x := range rawTags {
	//	fmt.Println(rawTags[x])
	//}

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
		return versionCompare(semverTags[i], semverTags[j])
	})

	// No tags found, start with initial version
	if len(semverTags) == 0 {
		fmt.Println("No tags found, creating the first tag v1.0.0-DEV.1")
		fmt.Println("::set-output name=new_tag::v1.0.0-DEV.1")
		os.Exit(0)
	}

	// Get the latest tag
	latestTag := semverTags[len(semverTags)-1]

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

// versionCompare compares two semantic versioning tags.
func versionCompare(tag1, tag2 string) bool {
	version1 := strings.TrimPrefix(tag1, "v")
	version2 := strings.TrimPrefix(tag2, "v")

	v1Parts := strings.Split(version1, ".")
	v2Parts := strings.Split(version2, ".")

	for i := 0; i < 3; i++ {
		v1, _ := strconv.Atoi(v1Parts[i])
		v2, _ := strconv.Atoi(v2Parts[i])
		if v1 != v2 {
			return v1 < v2
		}
	}

	return false
}
