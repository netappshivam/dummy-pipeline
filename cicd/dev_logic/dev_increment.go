package dev_logic

import (
	"dummy-pipeline/cicd/release_branch_logic"
	"fmt"
)

func Dev_increment_logic() {

	lastTag := release_branch_logic.LatestTagFetch("Dev")
	fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

	newTag, _ := release_branch_logic.IncrementTag(lastTag)
	if len(newTag) == 0 {
		newTag = "2520.0.0-DEV.1"
	}
	fmt.Printf("::set-output name=new_tag::%s\n", newTag)

}
