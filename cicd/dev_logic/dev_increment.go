package dev_logic

import (
	"dummy-pipeline/cicd/release_branch_logic"
	"fmt"
	"os"
)

func Dev_increment_logic() {

	err := release_branch_logic.FetchTags()
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	lastTag := release_branch_logic.LatestTagFetch("Dev")
	fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

	newTag, _ := release_branch_logic.IncrementTag(lastTag)
	if len(newTag) == 0 {
		newTag = fmt.Sprintf(release_branch_logic.SetupConfigobject.CurrentWeekRelease + ".0.0-DEV.1")
	}

	fmt.Printf("::set-output name=newTag1::%s\n", newTag)
}
