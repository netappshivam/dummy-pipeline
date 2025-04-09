package main

import "dummy-pipeline/release_branch_logic"

func main() {
	// Initialize data
	release_branch_logic.IntializeData()

	// Create a release branch
	release_branch_logic.Release_creation()
}
