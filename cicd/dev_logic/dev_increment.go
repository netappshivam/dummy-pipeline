package dev_logic

import (
	"dummy-pipeline/cicd/release_branch_logic"
	"fmt"
)

func Dev_increment_logic() {

	lastTag := release_branch_logic.LatestTagFetch("Dev")
	fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

	newTag, _ := release_branch_logic.IncrementTag(lastTag)
	fmt.Printf("New tag is: %s\n", newTag)

	release_branch_logic.SetNewTag(newTag)

}
