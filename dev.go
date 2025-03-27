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

	currentTagDateName := Xyz()

	// Fetch all tags from the remote
	cmd := exec.Command("git", "fetch", "--tags")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	// Get all tag names
	cmd = exec.Command("git", "tag", "--list", currentTagDateName+"*-DEV.*")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error listing tags:", err)
		os.Exit(1)
	}
	rawTags := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter tags that follow semantic versioning
	semverTags := make([]string, 0)
	for _, tag := range rawTags {
		if strings.HasPrefix(tag, currentTagDateName) {
			parts := strings.Split(strings.TrimPrefix(tag, currentTagDateName), "-")
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

		//tag := fmt.Sprintf("%02d%02d", currentYearStr, currentMonthStr)
		//fmt.Printf("No tags found, creating the first tag %s\n", tag+"0.0-DEV.1")
		//fmt.Printf("::set-output name=new_tag::%s\n", tag+"0.0-DEV.1")
		//os.Exit(0)

		fmt.Println("No tags found, creating the first tag" + currentTagDateName + ".0.0-DEV.1")
		fmt.Println("::set-output name=new_tag::" + currentTagDateName + ".0.0-DEV.1")
		os.Exit(0)
	}

	// Get the latest tag
	latestTag := semverTags[len(semverTags)-1]

	// Increment the tag
	parts := strings.Split(strings.TrimPrefix(latestTag, currentTagDateName), "-")
	if len(parts) != 2 {
		fmt.Println("Latest tag does not follow semantic versioning:", latestTag)
		os.Exit(1)
	}
	devParts := strings.Split(parts[1], ".")
	if len(devParts) != 2 || devParts[0] != "DEV" {
		fmt.Println("Latest tag does not follow semantic versioning:", latestTag)
		os.Exit(1)
	}
	devNumber, _ := strconv.Atoi(devParts[1])
	devNumber++
	newTag := fmt.Sprintf("v%s-DEV.%d", parts[0], devNumber)

	// Output the new tag to be used by other steps
	fmt.Printf("::set-output name=new_tag::%s\n", newTag)
}

// versionCompare_RC compares two semantic versioning tags.
func versionCompare(tag1, tag2 string) bool {
	currenTagDateName := Xyz()
	version1 := strings.TrimPrefix(tag1, currenTagDateName)
	version2 := strings.TrimPrefix(tag2, currenTagDateName)

	v1Parts := strings.Split(version1, "-")
	v2Parts := strings.Split(version2, "-")

	if v1Parts[0] != v2Parts[0] {
		return v1Parts[0] < v2Parts[0]
	}

	dev1, _ := strconv.Atoi(strings.TrimPrefix(v1Parts[1], "DEV."))
	dev2, _ := strconv.Atoi(strings.TrimPrefix(v2Parts[1], "DEV."))
	return dev1 < dev2
}
